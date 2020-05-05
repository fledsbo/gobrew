package hwinterface

import (
	"log"
	"sync"
	"time"

	"github.com/fledsbo/gatt"
	"github.com/fledsbo/gatt/examples/option"
)

// MonitorState represents the state of one monitor
type MonitorState struct {
	Name        string
	Type        string
	Timestamp   time.Time
	Gravity     *float64
	Temperature *float64
}

type monitorController interface {
	GetMonitor(string) (MonitorState, bool)
	GetMonitors() []MonitorState
	SetMonitor(MonitorState)
	Scan()
}

// The MonitorController maintains the state of all monitors
type MonitorController struct {
	mux      sync.Mutex
	c        chan MonitorState
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

	gravity := float64(b.minor) / 1000
	temperature := (float64(b.major) - 32) * 5 / 9

	m.SetMonitor(MonitorState{
		Name:        tilt,
		Type:        "Tilt",
		Timestamp:   time.Now(),
		Gravity:     &gravity,
		Temperature: &temperature,
	})
}

// NewMonitorController creates a new monitor controller
func NewMonitorController() (out *MonitorController) {
	out = new(MonitorController)
	out.Monitors = make(map[string]MonitorState)
	out.c = make(chan MonitorState, 100)
	return
}

// Scan runs forever, reading the state of each monitor
func (m *MonitorController) Scan() {
	log.Printf("Starting scan")
	device, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		panic(err)
	}
	device.Handle(gatt.PeripheralDiscovered(m.onPeripheralDiscovered))
	device.Init(onStateChanged)
	for ms := range m.c {
		m.mux.Lock()
		m.Monitors[ms.Name] = ms
		m.mux.Unlock()
	}
}

// GetMonitors return all monitor states
func (m *MonitorController) GetMonitors() (out []*MonitorState) {
	m.mux.Lock()
	defer m.mux.Unlock()
	out = make([]*MonitorState, 0, len(m.Monitors))
	for _, v := range m.Monitors {
		log.Printf("Returning data for %s", v.Name)
		m := v
		out = append(out, &m)
	}
	return
}

// GetMonitor returns a specific monitor state
func (m *MonitorController) GetMonitor(monitor string) (state MonitorState, found bool) {
	m.mux.Lock()
	defer m.mux.Unlock()
	state, found = m.Monitors[monitor]
	return
}

// SetMonitor sets the value of a specific monitor
func (m *MonitorController) SetMonitor(mstate MonitorState) {
	m.c <- mstate
}
