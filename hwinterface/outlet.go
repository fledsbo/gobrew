package hwinterface

import (
	"log"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/warthog618/gpiod"
)

// Outlet represents one outlet config
type Outlet struct {
	Name        string
	CodeOn      uint64
	CodeOff     uint64
	Protocol    uint
	PulseLength uint
}

// OutletController maintains a list of configured controllers and the interface to switch them on and off
type OutletController struct {
	chip        *gpiod.Chip
	offset      int
	transmitter *gpio.Transmitter
	Outlets     map[string]Outlet
}

func encode(str []byte) (out uint64) {
	for _, b := range str {
		out <<= 2
		switch b {
		case 'F':
			out |= 1
		case '0':
			out |= 0
		case '1':
			out |= 3
		}
	}
	return
}

// AddDialOutlet adds a dial-type outlet to the controller
func (c *OutletController) AddDialOutlet(name string, group int, id int) Outlet {
	outlet := Outlet{
		Name:        name,
		Protocol:    0,
		PulseLength: 350,
		CodeOn:      encode(dialCode(group, id, true)),
		CodeOff:     encode(dialCode(group, id, false)),
	}

	c.Outlets[name] = outlet
	return outlet
}

func (c *OutletController) GetOutlet(name string) *Outlet {
	outlet, found := c.Outlets[name]
	if found {
		return &outlet
	} else {
		return nil
	}
}

func dialCode(group int, id int, on bool) (out []byte) {
	out = []byte{'F', 'F', 'F', 'F', 'F', 'F', 'F', 'F', 'F', 'F', 'F', 'F'}
	out[group] = '0'
	out[id+4] = '0'
	if !on {
		out[11] = '0'
	}
	return
}

// NewOutletController creates a new OutletController
func NewOutletController() (out *OutletController) {
	out = new(OutletController)
	chip, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		panic(err)
	}
	out.chip = chip
	out.offset = 23

	transmitter, err := gpio.NewTransmitter(out.chip, out.offset, gpio.TransmissionCount(10))
	if err != nil {
		panic(err)
	}
	out.transmitter = transmitter

	out.Outlets = make(map[string]Outlet)
	return
}

func (t *OutletController) Close() {
	t.transmitter.Close()
}

func (t *OutletController) SwitchOn(outlet Outlet) {
	<-t.transmitter.Transmit(outlet.CodeOn, gpio.DefaultProtocols[outlet.Protocol], outlet.PulseLength)
}

func (t *OutletController) SwitchOff(outlet Outlet) {
	<-t.transmitter.Transmit(outlet.CodeOff, gpio.DefaultProtocols[outlet.Protocol], outlet.PulseLength)
}

func (t *OutletController) SwitchAllOff() {
	log.Println("Turning off all outlets")
	for _, outlet := range t.Outlets {
		t.SwitchOff(outlet)
	}
}
