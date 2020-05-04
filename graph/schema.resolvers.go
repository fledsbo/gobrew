package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/fledsbo/gobrew/graph/generated"
	"github.com/fledsbo/gobrew/graph/model"
)

func (r *queryResolver) Batches(ctx context.Context) ([]*model.Batch, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Monitors(ctx context.Context) ([]*model.FermentationMonitor, error) {
	mons := r.MonitorController.GetMonitors()
	outMons := make([]*model.FermentationMonitor, 0, len(mons))
	for _, m := range mons {
		outMons = append(outMons, &model.FermentationMonitor{m.Name, m.Name, "Tilt", &m.Temperature, &m.Gravity})
	}
	return outMons, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
