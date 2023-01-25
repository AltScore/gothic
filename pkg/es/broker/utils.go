package broker

import (
	"context"
	"fmt"
	"github.com/modernice/goes/event"
	"reflect"
)

// Subscribe register a given event handler in the provider broker to receive events with the given event name.
// The handler will receive events with the corresponding data type.
// If the event does not correspond to the expected data type, an error will be returned.
func Subscribe[T any](subscriber Broker, eventName string, handle TypedEventHandler[T]) (UnsubscribeFunc, error) {
	hf := func(ctx context.Context, ev event.Event) error {
		typedEvent, ok := event.TryCast[T, any](ev)

		if !ok {
			return fmt.Errorf("event %s is not of type %T", ev.Name(), GetTypeName(typedEvent))
		}

		return handle(ctx, typedEvent)
	}

	return subscriber.Subscribe(hf, eventName)
}

func GetTypeName(aVar any) string {
	if t := reflect.TypeOf(aVar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
