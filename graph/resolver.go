package graph

import (
	"github.com/fledsbo/gobrew/config"
	"github.com/fledsbo/gobrew/fermentation"
	"github.com/fledsbo/gobrew/hwinterface"
	"github.com/fledsbo/gobrew/storage"
)

// This file will not be regenerated automatically.

// Resolver serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct {
	Config                 *config.Config
	MonitorController      *hwinterface.MonitorController
	OutletController       *hwinterface.OutletController
	FermentationController *fermentation.Controller
	Storage                *storage.Storage
}
