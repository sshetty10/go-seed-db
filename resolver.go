package main

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

//Query returns the QueryResolver interface
//QueryResolver is an interface in generated.go

/*func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct {
	*Resolver
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
}*/
