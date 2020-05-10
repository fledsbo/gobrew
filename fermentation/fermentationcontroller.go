package fermentation

import (
	"log"
	"time"

	"github.com/fledsbo/gobrew/hwinterface"
)

type Batch struct {
	Name string

	AssignedCoolingOutlet string
	AssignedHeatingOutlet string
	AssignedMonitor       string

	TargetTemp        float64
	Hysteresis        float64
	MaxReadingAge     time.Duration
	MinOutletDuration time.Duration

	CurrentlyHeating bool
	CurrentlyCooling bool
	LastStateChange  time.Time

	CurrentGravity float64
	CurrentTemp    float64

	tempAverager    averager
	gravityAverager averager
}

type BatchState struct {
	Temperature *float64
	Gravity     *float64
}

type Controller struct {
	Batches []*Batch

	monitorController *hwinterface.MonitorController
	outletController  *hwinterface.OutletController
}

func NewController(monitorC *hwinterface.MonitorController, outletC *hwinterface.OutletController) (out *Controller) {
	out = &Controller{
		monitorController: monitorC,
		outletController:  outletC,
	}

	go out.Run()
	return
}

func (c *Controller) GetBatchState(batch *Batch) (out BatchState) {

	return BatchState{
		&batch.CurrentTemp,
		&batch.CurrentGravity,
	}
}

func (c *Batch) check(monitorController *hwinterface.MonitorController, outletController *hwinterface.OutletController) {
	var monitorState hwinterface.MonitorState
	found := false

	if c.AssignedMonitor != "" {
		monitorState, found = monitorController.GetMonitor(c.AssignedMonitor)
	}

	previousCooling := c.CurrentlyCooling
	previousHeating := c.CurrentlyHeating

	if found && monitorState.Temperature != nil && time.Now().Sub(monitorState.Timestamp) < c.MaxReadingAge {

		c.CurrentTemp = c.tempAverager.addReading(*monitorState.Temperature)

		if monitorState.Gravity != nil {
			c.CurrentGravity = c.gravityAverager.addReading(*monitorState.Gravity)
		}

		if c.CurrentTemp > c.TargetTemp &&
			time.Now().Sub(c.LastStateChange) > c.MinOutletDuration {
			c.CurrentlyHeating = false
			if c.CurrentTemp > (c.TargetTemp + c.Hysteresis) {
				c.CurrentlyCooling = true
			}
			c.LastStateChange = time.Now()
		}
		if c.CurrentTemp < c.TargetTemp &&
			time.Now().Sub(c.LastStateChange) > c.MinOutletDuration {
			c.CurrentlyCooling = false
			if c.CurrentTemp < (c.TargetTemp - c.Hysteresis) {
				c.CurrentlyHeating = true
			}
			c.LastStateChange = time.Now()
		}

	} else {
		// We don't have a monitor, or we're not getting readings.
		// Turn everything off to be safe
		c.CurrentlyHeating = false
		c.CurrentlyCooling = false
	}

	if previousCooling != c.CurrentlyCooling {
		if c.CurrentlyCooling {
			log.Printf("Turning cooling on")
		} else {
			log.Printf("Turning cooling off")
		}
	}
	if previousHeating != c.CurrentlyHeating {
		if c.CurrentlyHeating {
			log.Printf("Turning heating on")
		} else {
			log.Printf("Turning heating off")
		}
	}

	heatingOutlet := outletController.GetOutlet(c.AssignedHeatingOutlet)

	// We set the outlets every time, in case they missed a previous command
	if heatingOutlet != nil {
		if c.CurrentlyHeating {
			outletController.SwitchOn(*heatingOutlet)
		} else {
			outletController.SwitchOff(*heatingOutlet)
		}
	}

	coolingOutlet := outletController.GetOutlet(c.AssignedCoolingOutlet)

	if coolingOutlet != nil {
		if c.CurrentlyCooling {
			outletController.SwitchOn(*coolingOutlet)
		} else {
			outletController.SwitchOff(*coolingOutlet)
		}
	}
}

func (c *Controller) Run() {
	for {
		for _, b := range c.Batches {
			b.check(c.monitorController, c.outletController)
		}
		time.Sleep(time.Second)
	}
}
