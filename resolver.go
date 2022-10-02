package main

import (
	"context"
	"strings"

	"github.com/sshetty10/go-seed-db/generated"
	"github.com/sshetty10/go-seed-db/model"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	api *API
}

// NewResolver makes a new Resolver
//
func (api *API) NewResolver() *Resolver {
	return &Resolver{
		api: api,
	}
}

//needed a separate GQL trainer cause the ID field is primitive.ObjectID for Trainer
//without that type we cannot get the ID back from DB (bson mapping for ID field needs it to be of the primitive.ObjectID type)
func (r *queryResolver) Trainers(ctx context.Context) ([]*model.Trainer, error) {
	trainers, err := r.api.ListDBTrainers()
	if err != nil {
		return nil, err
	}

	return trainers, nil
}
func (r *queryResolver) TrainerByID(ctx context.Context, id string) (*model.Trainer, error) {
	t, err := r.api.GetDBTrainerByID(id)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *trainerResolver) LicenseState(ctx context.Context, trainer *model.Trainer) (string, error) {
	licenseState := strings.Split(trainer.LicenseID, "-")[0]

	//pprofed

	/*idx := 0
	for idx < len(trainer.LicenseID) {
		if trainer.LicenseID[idx] == '-' {
			break
		}
		idx++
	}
	licenseState := trainer.LicenseID[:idx]*/

	return licenseState, nil
}

//Query returns the QueryResolver interface
//QueryResolver is an interface in generated.go

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

// Trainer returns generated.TrainerResolver implementation.
func (r *Resolver) Trainer() generated.TrainerResolver { return &trainerResolver{r} }

type queryResolver struct{ *Resolver }
type trainerResolver struct{ *Resolver }
