package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

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

func (r *queryResolver) Batches(ctx context.Context) ([]*model.Batch, error) {
	return []*model.Batch{}, nil
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *fermentationMonitorResolver) ID(ctx context.Context, obj *hwinterface.MonitorState) (string, error) {
	return obj.Name, nil
}
