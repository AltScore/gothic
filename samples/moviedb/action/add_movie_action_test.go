package action

import (
	"context"
	"testing"

	"github.com/AltScore/gothic/pkg/actions/optimistic"
	"github.com/AltScore/gothic/pkg/apiusers"
	"github.com/AltScore/gothic/pkg/apiusers/basicuser"
	"github.com/AltScore/gothic/samples/moviedb/model"
	"github.com/stretchr/testify/assert"
)

var testUser = basicuser.New("test", "test", "test", []string{})

func TestAddMovieAction_Execute(t *testing.T) {
	action := NewAddMovieAction()

	execute, err := action.Execute(givenTestContext(), model.AddMovieRequest{})

	assert.NoError(t, err)
	assert.Equal(t, model.AddMovieResponse{}, execute)
}

func TestAddMovieAction_ExecuteOptimistically(t *testing.T) {
	// Cannot use 'optimistic.Wrapp[model.AddMovieRequest, model.AddMovieResponse]'
	//       (type func[Request any, Response any](target actions.Action[Request, Response]) *Wrapper[Request, Response])
	// as the type func(Action[Request, Response]) Action[Request, Response]
	action := optimistic.Wrapp[model.AddMovieRequest, model.AddMovieResponse](NewAddMovieAction())

	execute, err := action.Execute(givenTestContext(), model.AddMovieRequest{})

	assert.NoError(t, err)
	assert.Equal(t, model.AddMovieResponse{}, execute)
}

func givenTestContext() context.Context {
	return apiusers.SetApiUser(context.Background(), testUser)
}
