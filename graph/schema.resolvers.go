package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fledsbo/gobrew/fermentation"
	"github.com/fledsbo/gobrew/graph/generated"
	"github.com/fledsbo/gobrew/graph/model"
	"github.com/fledsbo/gobrew/hwinterface"
)

func (r *fermentationMonitorResolver) Timestamp(ctx context.Context, obj *hwinterface.MonitorState) (*string, error) {
	timestamp := obj.Timestamp.Format("2006-01-02T15:04:05-0700")
	return &timestamp, nil
}

func (r *mutationResolver) SetMonitor(ctx context.Context, input *model.SetMonitorInput) (string, error) {
	r.MonitorController.SetMonitor(hwinterface.MonitorState{
		Name:        input.Name,
		Timestamp:   time.Now(),
		Gravity:     input.Gravity,
		Temperature: input.Temperature,
	})
	return input.Name, nil
}

func (r *mutationResolver) SetFermentation(ctx context.Context, input *model.SetFermentationInput) (string, error) {
	var found *fermentation.Batch
	for _, ferm := range r.FermentationController.Batches {
		if ferm.Name == input.Name {
			found = ferm
		} else {
			if input.Monitor != nil && ferm.AssignedMonitor == *input.Monitor {
				return "", fmt.Errorf("Monitor %s already used by %s", *input.Monitor, ferm.Name)
			}

			if input.HeatingOutlet != nil && ferm.AssignedCoolingOutlet == *input.HeatingOutlet || ferm.AssignedHeatingOutlet == *input.HeatingOutlet {
				return "", fmt.Errorf("Outlet %s already used by %s", *input.HeatingOutlet, ferm.Name)
			}

			if input.CoolingOutlet != nil && ferm.AssignedCoolingOutlet == *input.CoolingOutlet || ferm.AssignedHeatingOutlet == *input.CoolingOutlet {
				return "", fmt.Errorf("Outlet %s already used by %s", *input.CoolingOutlet, ferm.Name)
			}
		}
	}

	if found == nil {
		found = &fermentation.Batch{
			Name: input.Name,
		}
		r.FermentationController.Batches = append(r.FermentationController.Batches, found)
	}

	if input.Monitor != nil {
		found.AssignedMonitor = *input.Monitor
	}

	if input.HeatingOutlet != nil {
		if found.AssignedCoolingOutlet == *input.HeatingOutlet {
			return "", fmt.Errorf("Cannot use same outlet %s for heating and cooling", *input.HeatingOutlet)
		}
		found.AssignedHeatingOutlet = *input.HeatingOutlet
	}

	if input.CoolingOutlet != nil {
		if found.AssignedHeatingOutlet == *input.CoolingOutlet {
			return "", fmt.Errorf("Cannot use same outlet %s for heating and cooling", *input.CoolingOutlet)
		}
		found.AssignedCoolingOutlet = *input.CoolingOutlet
	}

	if input.Config != nil {
		if input.Config.TargetTemp != nil {
			found.TargetTemp = *input.Config.TargetTemp
		}
		if input.Config.Hysteresis != nil {
			found.Hysteresis = *input.Config.Hysteresis
		}
		if input.Config.MaxReadingAgeSec != nil {
			found.MaxReadingAge = time.Duration(*input.Config.MaxReadingAgeSec) * time.Second
		}
		if input.Config.MinOutletDurationSec != nil {
			found.MinOutletDuration = time.Duration(*input.Config.MinOutletDurationSec) * time.Second
		}
	}

	err := r.Storage.StoreFermentations()

	return found.Name, err
}

func (r *mutationResolver) RemoveFermentation(ctx context.Context, input *model.RemoveFermentationInput) (string, error) {
	for i, b := range r.FermentationController.Batches {
		if b.Name == input.Name {
			r.FermentationController.Batches = append(r.FermentationController.Batches[:i], r.FermentationController.Batches[i+1:]...)
			r.Storage.RemoveFermentation(input.Name)
			return input.Name, nil
		}
	}
	return "", errors.New("Fermentation not found")
}

func (r *mutationResolver) SetupDialOutlet(ctx context.Context, input *model.SetupDialOutletInput) (string, error) {
	r.OutletController.AddDialOutlet(input.Name, input.Group, input.Outlet)
	err := r.Storage.StoreOutlets()
	return input.Name, err
}

func (r *queryResolver) Fermentations(ctx context.Context) ([]*model.Fermentation, error) {
	out := make([]*model.Fermentation, 0, len(r.FermentationController.Batches))
	for _, fc := range r.FermentationController.Batches {
		f := model.Fermentation{
			Name:    fc.Name,
			CanHeat: fc.AssignedHeatingOutlet != "",
			Heating: fc.CurrentlyHeating,
			CanCool: fc.AssignedCoolingOutlet != "",
			Cooling: fc.CurrentlyCooling,
		}
		config := model.FermentationConfig{
			TargetTemp:           fc.TargetTemp,
			Hysteresis:           fc.Hysteresis,
			MaxReadingAgeSec:     int(fc.MaxReadingAge.Seconds()),
			MinOutletDurationSec: int(fc.MinOutletDuration.Seconds()),
		}
		f.Config = &config
		monitor, found := r.MonitorController.GetMonitor(fc.AssignedMonitor)
		if found {
			f.Monitor = &monitor
		}
		out = append(out, &f)
	}
	return out, nil
}

func (r *queryResolver) Monitors(ctx context.Context) ([]*hwinterface.MonitorState, error) {
	mons := r.MonitorController.GetMonitors()
	return mons, nil
}

func (r *queryResolver) Outlets(ctx context.Context) ([]*hwinterface.Outlet, error) {
	out := make([]*hwinterface.Outlet, 0, len(r.OutletController.Outlets))
	for _, outlet := range r.OutletController.Outlets {
		o := outlet
		out = append(out, &o)
	}
	return out, nil
}

// FermentationMonitor returns generated.FermentationMonitorResolver implementation.
func (r *Resolver) FermentationMonitor() generated.FermentationMonitorResolver {
	return &fermentationMonitorResolver{r}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type fermentationMonitorResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
