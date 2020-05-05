package graph

import (
	"github.com/fledsbo/gobrew/fermentation"
	"github.com/fledsbo/gobrew/hwinterface"
	"github.com/fledsbo/gobrew/storage"
)

// This file will not be regenerated automatically.

// Resolver serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct {
	MonitorController       *hwinterface.MonitorController
	OutletController        *hwinterface.OutletController
	FermentationControllers []*fermentation.FermentationController
	Storage                 *storage.Storage
}
