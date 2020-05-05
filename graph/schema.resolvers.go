package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"
	"errors"

	"github.com/fledsbo/gobrew/graph/generated"
	"github.com/fledsbo/gobrew/graph/model"
	"github.com/fledsbo/gobrew/hwinterface"
	"github.com/fledsbo/gobrew/fermentation"
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
	var found *fermentation.FermentationController
	for _,ferm := range r.FermentationControllers {
		if ferm.Name == input.Name {
			found = ferm			
		} else {
			if (input.Monitor != nil && ferm.AssignedMonitor == *input.Monitor) {
				return "", errors.New(fmt.Sprintf("Monitor %s already used by %s", *input.Monitor, ferm.Name))
			}						
		}
	}

	if found == nil {
		found = fermentation.NewFermentationController(input.Name, r.MonitorController, r.OutletController)
		r.FermentationControllers = append(r.FermentationControllers, found)
	}

	if input.Monitor != nil {
		found.AssignedMonitor = *input.Monitor
	}

	if input.Config != nil {
		if (input.Config.TargetTemp != nil) {
			found.TargetTemp = *input.Config.TargetTemp
		}
		if (input.Config.Hysteresis != nil) {
			found.Hysteresis = *input.Config.Hysteresis
		}
		if (input.Config.MaxReadingAgeSec != nil) {
			found.MaxReadingAge = time.Duration(*input.Config.MaxReadingAgeSec) * time.Second
		}
		if (input.Config.MinOutletDurationSec != nil) {
			found.MinOutletDuration = time.Duration(*input.Config.MinOutletDurationSec) * time.Second
		}
	}

	return found.Name, nil
}

func (r *queryResolver) Fermentations(ctx context.Context) ([]*model.Fermentation, error) {
	out := make([]*model.Fermentation, 0, len(r.FermentationControllers))
	for _, fc := range r.FermentationControllers {
		f := model.Fermentation {
			Name: fc.Name,
			CanHeat: fc.HeatingOutlet != nil,
			Heating: fc.CurrentlyHeating,
			CanCool: fc.CoolingOutlet != nil,
			Cooling: fc.CurrentlyCooling,
		}
		config := model.FermentationConfig {
			TargetTemp: fc.TargetTemp,
			Hysteresis: fc.Hysteresis,
			MaxReadingAgeSec: int(fc.MaxReadingAge.Seconds()),
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

