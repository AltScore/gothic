package broker

import (
	"context"
	"github.com/modernice/goes/event"
)

// UnsubscribeFunc is a function that when called cancel the original subscription
type UnsubscribeFunc func()
type TypedEventHandler[T any] func(ctx context.Context, event event.Of[T]) error
