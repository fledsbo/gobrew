package fermentation

import (
	"log"
	"time"

	"github.com/fledsbo/gobrew/hwinterface"
)

type FermentationController struct {
	Name string

	MonitorController *hwinterface.MonitorController
	OutletController  *hwinterface.OutletController

	HeatingOutlet   *hwinterface.Outlet
	CoolingOutlet   *hwinterface.Outlet

	AssignedMonitor 	  string

	TargetTemp        float64
	Hysteresis        float64
	MaxReadingAge     time.Duration
	MinOutletDuration time.Duration

	CurrentGravity float64

	CurrentlyHeating bool
	CurrentlyCooling bool
	LastStateChange  time.Time
}

func NewFermentationController(name string, monitorC *hwinterface.MonitorController, outletC *hwinterface.OutletController) (out *FermentationController) {
	out = &FermentationController{
		Name: name,

		MonitorController: monitorC,
		OutletController:  outletC,

		TargetTemp: 18.0,
		Hysteresis: 0.5,

		MaxReadingAge:     10 * time.Minute,
		MinOutletDuration: 1 * time.Minute,

		CurrentlyHeating: false,
		CurrentlyCooling: false,
	}

	heatOutlet := hwinterface.NewDialOutlet(0, 0)
	coolOutlet := hwinterface.NewDialOutlet(0, 1)
	
	out.HeatingOutlet = &heatOutlet
	out.CoolingOutlet = &coolOutlet


	go out.Run()
	return
}

func (c *FermentationController) Check() {
	var monitorState hwinterface.MonitorState
	found := false

	if c.AssignedMonitor != "" {
		monitorState, found = c.MonitorController.GetMonitor(c.AssignedMonitor)
	}

	previousCooling := c.CurrentlyCooling
	previousHeating := c.CurrentlyHeating

	if found && monitorState.Temperature != nil && time.Now().Sub(monitorState.Timestamp) < c.MaxReadingAge {
		if *monitorState.Temperature > c.TargetTemp &&
			time.Now().Sub(c.LastStateChange) > c.MinOutletDuration {
			c.CurrentlyHeating = false
			if *monitorState.Temperature > (c.TargetTemp + c.Hysteresis) {
				c.CurrentlyCooling = true
			}
			c.LastStateChange = time.Now()
		}
		if *monitorState.Temperature < c.TargetTemp &&
			time.Now().Sub(c.LastStateChange) > c.MinOutletDuration {
			c.CurrentlyCooling = false
			if *monitorState.Temperature < (c.TargetTemp - c.Hysteresis) {
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

	// We set the outlets every time, in case they missed a previous command
	if c.HeatingOutlet != nil {
		if c.CurrentlyHeating {
			c.OutletController.SwitchOn(*c.HeatingOutlet)
		} else {
			c.OutletController.SwitchOff(*c.HeatingOutlet)
		}
	}

	if c.CoolingOutlet != nil {
		if c.CurrentlyCooling {
			c.OutletController.SwitchOn(*c.CoolingOutlet)
		} else {
			c.OutletController.SwitchOff(*c.CoolingOutlet)
		}
	}
}

func (c *FermentationController) Run() {
	for {
		c.Check()
		time.Sleep(time.Second)
	}
}
