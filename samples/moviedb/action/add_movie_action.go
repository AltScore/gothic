package action

import (
	"context"

	"github.com/AltScore/gothic/pkg/actions"
	"github.com/AltScore/gothic/samples/moviedb/model"
)

type AddMovieAction struct {
	actions.Action[model.AddMovieRequest, model.AddMovieResponse]
}

func NewAddMovieAction() *AddMovieAction {
	return &AddMovieAction{
		actions.New[model.AddMovieRequest, model.AddMovieResponse](),
	}
}

func (a AddMovieAction) Execute(ctx context.Context, request model.AddMovieRequest) (model.AddMovieResponse, error) {
	user := a.User(ctx)
	a.Infof("AddMovieAction.Execute %s", user)
	return model.AddMovieResponse{}, nil
}
