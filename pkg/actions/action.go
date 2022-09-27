package actions

import (
	"context"

	"github.com/AltScore/gothic/pkg/apiusers"
	"github.com/AltScore/gothic/pkg/loggers"
)

type Wrapper[Request any, Response any] func(target Action[Request, Response]) Action[Request, Response]

type Action[Request any, Response any] interface {
	loggers.Logger
	Execute(context.Context, Request) (Response, error)
	User(ctx context.Context) apiusers.ApiUser
}

type RootAction[Request any, Response any] struct {
	loggers.Logger
}

func New[Request any, Response any]() RootAction[Request, Response] {
	return RootAction[Request, Response]{
		Logger: loggers.NewSimple(),
	}
}

func (a RootAction[Request, Response]) Execute(_ context.Context, _ Request) (Response, error) {
	panic("should be overridden")
}

func (a RootAction[Request, Response]) User(ctx context.Context) apiusers.ApiUser {
	return apiusers.GetApiUser(ctx)
}
