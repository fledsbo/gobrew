package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
	"time"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/warthog618/gpiod"
)

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

func scanTilt() {
	log.SetOutput(ioutil.Discard)
	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
	log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	device.Handle(gatt.PeripheralDiscovered(onPeripheralDiscovered))
	device.Init(onStateChanged)
	select {}
}

func testOutlets() {
	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		panic(err)
	}

	offset := 23
	transmitter, err := gpio.NewTransmitter(c, offset, gpio.TransmissionCount(10))
	if err != nil {
		panic(err)
	}

	defer transmitter.Close()


	for g := 0; g < 4; g++ {
		for o := 0; o < 4; o++ {
			fmt.Printf("Testing %d %d\n", g, o)
			ol := NewDialOutlet(g, o)
			<-transmitter.Transmit(ol.CodeOn, gpio.DefaultProtocols[ol.Protocol], ol.PulseLength)
			time.Sleep(4 * time.Second)
			<-transmitter.Transmit(ol.CodeOff, gpio.DefaultProtocols[ol.Protocol], ol.PulseLength)
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	testOutlets()
}

