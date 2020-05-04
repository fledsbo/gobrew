package main

import (
	"fmt"
	//"github.com/paypal/gatt"
)

/*
func onStateChanged(device gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning for iBeacon Broadcasts...")
		device.Scan([]gatt.UUID{}, true)
		return
	default:
		device.StopScanning()
	}
}

func onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	tiltIds := map[string]string{
		"A495BB10-C5B1-4B44-B512-1370F02D74DE": "Red",
		"A495BB20-C5B1-4B44-B512-1370F02D74DE": "Green",
		"A495BB30-C5B1-4B44-B512-1370F02D74DE": "Black",
		"A495BB40-C5B1-4B44-B512-1370F02D74DE": "Purple",
		"A495BB50-C5B1-4B44-B512-1370F02D74DE": "Orange",
		"A495BB60-C5B1-4B44-B512-1370F02D74DE": "Blue",
		"A495BB70-C5B1-4B44-B512-1370F02D74DE": "Yellow",
		"A495BB80-C5B1-4B44-B512-1370F02D74DE": "Pink",
	}
	b, err := NewiBeacon(a.ManufacturerData)
	if err != nil {
		return
	}

	tilt, ok := tiltIds[b.uuid]
	if !ok {
		return
	}

	fmt.Printf("%s: %.1f %d\r", tilt, (float64(b.major)-32)*5/9, b.minor)
}
*/

func main() {
	//log.SetOutput(ioutil.Discard)
	//device, err := gatt.NewDevice(option.DefaultClientOptions...)
	//if err != nil {
	//log.Fatalf("Failed to open device, err: %s\n", err)
	//return
	//}
	//device.Handle(gatt.PeripheralDiscovered(onPeripheralDiscovered))
	//device.Init(onStateChanged)
	//select {}

	//c, err := gpiod.NewChip("gpiochip0")
	//if err != nil {
	//	panic(err)
	//}

	//offset := rpi.GPIO4
	//transmitter := gpio.NewTransmitter(c, offset, gpio.TransmissionCount(1))
	//defer transmitter.Close()
	//
	//	<-transmitter.Transmit(

	ol1 := NewDialOutlet(0, 0)
	ol2 := NewDialOutlet(0, 1)

	fmt.Println(ol1.CodeOn, ol1.CodeOff)
	fmt.Println(ol2.CodeOn, ol2.CodeOff)
}
