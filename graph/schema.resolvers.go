package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/fledsbo/gobrew/graph/generated"
	"github.com/fledsbo/gobrew/graph/model"
	"github.com/fledsbo/gobrew/hwinterface"
)

func (r *fermentationMonitorResolver) ID(ctx context.Context, obj *hwinterface.MonitorState) (string, error) {
	return obj.Name, nil
}

func (r *fermentationMonitorResolver) Timestamp(ctx context.Context, obj *hwinterface.MonitorState) (*string, error) {
	timestamp := obj.Timestamp.Format("2006-01-02T15:04:05-0700")
	return &timestamp, nil
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

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type fermentationMonitorResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
