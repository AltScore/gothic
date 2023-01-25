package bus

import (
	"context"
	"github.com/modernice/goes/event"
)

type EventHandler func(ctx context.Context, event event.Event) error
