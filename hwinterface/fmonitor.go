package hwinterface

import (
	"sync"
	"time"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

type MonitorState struct {
	Name        string
	LastRead    time.Time
	Gravity     float64
	Temperature float64
}

type MonitorController struct {
	mux      sync.Mutex
	Monitors map[string]MonitorState
}

func onStateChanged(device gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		device.Scan([]gatt.UUID{}, true)
		return
	default:
		device.StopScanning()
	}
}

func (m *MonitorController) onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
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

	m.mux.Lock()
	defer m.mux.Unlock()
	m.Monitors[tilt] = MonitorState{tilt, time.Now(), float64(b.minor) / 1000, (float64(b.major) - 32) * 5 / 9}
}

func (m *MonitorController) Scan() {
	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		panic(err)
	}
	device.Handle(gatt.PeripheralDiscovered(m.onPeripheralDiscovered))
	device.Init(onStateChanged)
	select {}
}

func (m *MonitorController) GetMonitors() (out []MonitorState) {
	m.mux.Lock()
	defer m.mux.Unlock()
	out = make([]MonitorState, 0, len(m.Monitors))
	for _, v := range m.Monitors {
		out = append(out, v)
	}
	return
}
