package optimistic

import (
	"context"

	"github.com/AltScore/gothic/pkg/actions"
)

type Wrapper[Request any, Response any] struct {
	actions.Action[Request, Response]
}

func Wrapp[Request any, Response any](target actions.Action[Request, Response]) actions.Action[Request, Response] {
	return &Wrapper[Request, Response]{target}
}

func (a *Wrapper[Request, Response]) Execute(ctx context.Context, request Request) (Response, error) {
	user := a.Action.User(ctx)

	a.Infof("ExecuteOptimistically start for %s", user)

	execute, err := a.Action.Execute(ctx, request)

	a.Infof("ExecuteOptimistically end")

	return execute, err
}
