package broker

import (
	"context"
	"fmt"
	"github.com/AltScore/gothic/pkg/es/bus"
	"github.com/modernice/goes/event"
	"github.com/modernice/goes/helper/pick"
)

// A Registerer can register handlers for different events. Types that implement
// Registerer can be passed to RegisterHandler(), ApplyWith(), and HandleWith()
// to conveniently register handlers for events.
//
//	var reg event.Registerer
//	event.RegisterEventHandler(reg, "foo", func(e event.Of[FooEvent]) {
//		log.Printf("handled %q event with data %v", e.Name(), e.Data())
//	}
//
// ApplyWith() and HandleWith() are aliases for RegisterHandler(), to allow for
// more concise code.
type Registerer interface {
	// RegisterEventHandler registers an event handler for the given event name.
	RegisterEventHandler(eventName string, handler bus.EventHandler)
}

// Handlers is a map of event names to event handlers. Handlers can be embedded
// into structs to implement [Registerer]. [*github.com/modernice/goes/aggregate.Base]
// embeds Handlers to allow for convenient registration of event handlers.
type Handlers map[string][]bus.EventHandler

// RegisterEventHandler implements [Registerer].
func (h Handlers) RegisterEventHandler(eventName string, handler bus.EventHandler) {
	h[eventName] = append(h[eventName], handler)
}

// EventHandlers returns the handlers for the given event.
func (h Handlers) EventHandlers(eventName string) []bus.EventHandler {
	return h[eventName]
}

// HandleEvent calls each registered handler of the given [Event].
func (h Handlers) HandleEvent(ctx context.Context, evt event.Event) error {
	handlers := h.EventHandlers(evt.Name())
	for _, handler := range handlers {
		if err := handler(ctx, evt); err != nil {
			return err
		}
	}
	return nil
}

// RegisterHandler registers the handler for the given event.
//
// Example using *aggregate.Base:
//
//	type Foo struct {
//		*aggregate.Base
//
//		Foo string
//		Bar int
//		Baz bool
//	}
//
//	type FooEvent { Foo string }
//	type BazEvent { Baz bool }
//
//	func NewFoo(id uuid.UUID) *Foo  {
//		foo := &Foo{Base: aggregate.New("foo", id)}
//
//		event.RegisterHandler(foo, "foo", foo.applyFoo)
//		event.RegisterHandler(foo, "bar", foo.applyBar)
//		event.RegisterHandler(foo, "baz", foo.applyBaz)
//
//		return foo
//	}
//
//	func (f *Foo) applyFoo(e event.Of[FooEvent]) {
//		f.Foo = e.Data().Foo
//	}
//
//	func (f *Foo) applyBar(e event.Of[int]) {
//		f.Bar = e.Data()
//	}
//
//	func (f *Foo) applyBaz(e event.Of[BazEvent]) {
//		f.Baz = e.Data().Baz
//	}
func RegisterHandler[Data any](r Registerer, eventName string, handler TypedEventHandler[Data]) {
	r.RegisterEventHandler(eventName, func(ctx context.Context, evt event.Event) error {
		if casted, ok := event.TryCast[Data](evt); ok {
			return handler(ctx, casted)
		}

		aggregateName := "<unknown>"
		if a, ok := r.(pick.AggregateProvider); ok {
			aggregateName = pick.AggregateName(a)
		}
		var zero Data
		panic(fmt.Errorf(
			"[goes/event.RegisterHandler] Cannot cast %T to %T. "+
				"You probably provided the wrong event name for this handler. "+
				"[event=%v, aggregate=%v]",
			evt.Data(), zero, eventName, aggregateName,
		))
	})
}

// ApplyWith is an alias for RegisterHandler. Use ApplyWith instead of
// RegisterHandler to make code more concise:
//
//	type Foo struct {
//		*projection.Base
//
//		Foo string
//	}
//
//	func NewFoo() *Foo  {
//		foo := &Foo{Base: projection.New()}
//
//		// Because we "apply" events to the projection.
//		event.ApplyWith(foo, foo.applyFoo, "foo", "buzz")
//
//		return foo
//	}
//
//	func (f *Foo) applyFoo(e event.Of[string]) {
//		f.Foo = e.Data()
//	}
func ApplyWith[Data any](r Registerer, handler TypedEventHandler[Data], eventNames ...string) {
	for _, name := range eventNames {
		RegisterHandler(r, name, handler)
	}
}

// HandleWith is an alias for RegisterHandler. Use HandleWith instead of
// RegisterHandler to make code more concise:
//
//	import "github.com/modernice/goes/event/handler"
//
//	var bus event.Bus
//	h := handler.New(bus)
//
//	event.HandleWith(h, h.handleFoo, "foo")
func HandleWith[Data any](r Registerer, handler TypedEventHandler[Data], eventNames ...string) {
	ApplyWith(r, handler, eventNames...)
}
