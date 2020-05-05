package hwinterface

import (
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/warthog618/gpiod"
)

type Outlet struct {
	CodeOn      uint64
	CodeOff     uint64
	Protocol    uint
	PulseLength uint
}

type OutletController struct {
	chip        *gpiod.Chip
	offset      int
	transmitter *gpio.Transmitter
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

func NewDialOutlet(group int, id int) Outlet {
	return Outlet{
		Protocol:    0,
		PulseLength: 350,
		CodeOn:      encode(dialCode(group, id, true)),
		CodeOff:     encode(dialCode(group, id, false)),
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

func NewOutletController() (out *OutletController) {
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
