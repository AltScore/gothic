package handler

import (
	"github.com/AltScore/gothic/pkg/handlers"
	"github.com/AltScore/gothic/samples/moviedb/model"
)

type AddMovieHandler struct {
	handlers.Handler[model.AddMovieRequest, model.AddMovieResponse]
}

func NewAddMovieHandler() *AddMovieHandler {
	return &AddMovieHandler{
		handlers.New(handlers.HandlerSpec{
			Name:        "AddMovie",
			Route:       "/movies",
			Method:      "POST",
			Description: "Add a new movie to the database",
		}),
	}
}
