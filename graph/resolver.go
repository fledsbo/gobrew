package graph

import (
	"github.com/fledsbo/gobrew/hwinterface"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	MonitorController *hwinterface.MonitorController
}
