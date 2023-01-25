package broker

import (
	"context"
	"fmt"
	"github.com/AltScore/gothic/pkg/es/bus"
	utils "github.com/AltScore/gothic/pkg/logger"
	"github.com/modernice/goes/event"
	"github.com/totemcaf/gollections/slices"
	"github.com/totemcaf/gollections/syncs"
	"go.uber.org/zap"
	"sync"
)

type Broker struct {
	logger      utils.Logger
	lock        *sync.RWMutex
	handlersMap Handlers
}

func New(logger utils.Logger) *Broker {
	return &Broker{
		logger:      logger,
		lock:        &sync.RWMutex{},
		handlersMap: Handlers{},
	}
}

func doNothing() {
	// do nothing
}

// Subscribe a handler to receive events with the given event name
func (b *Broker) Subscribe(handler bus.EventHandler, names ...string) (UnsubscribeFunc, error) {
	eventNameSet := slicesToSet(names)

	if len(eventNameSet) == 0 {
		return doNothing, nil
	}

	if len(eventNameSet) != len(names) {
		return doNothing, fmt.Errorf("duplicate event names: %s", names)
	}

	unsubscribeFuncs := make([]UnsubscribeFunc, 0, len(names))
	defer b.unsubscribeAll(unsubscribeFuncs)

	for _, name := range names {
		unsubscribe, err := b.subscribeToEvent(name, handler)

		if err != nil {
			return doNothing, err
		}

		unsubscribeFuncs = append(unsubscribeFuncs, unsubscribe)
	}

	if len(unsubscribeFuncs) == 1 {
		// Optimization in case there is only one subscription to rollback
		return unsubscribeFuncs[0], nil
	}

	unsubscribeAll := func() {
		b.unsubscribeAll(unsubscribeFuncs)
	}
	return unsubscribeAll, nil
}

func (b *Broker) subscribeToEvent(eventName string, handler bus.EventHandler) (UnsubscribeFunc, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.handlersMap[eventName] = append(b.handlersMap[eventName], handler)

	unsubscribe := func() {
		// b.unsubscribe(eventName, handler) // TODO solve how to allow unsubscribe for a handler, we cannot compare functions
	}

	return unsubscribe, nil
}

// unsubscribe removes the given handler from the list of handlers in the named subscription
func (b *Broker) unsubscribe(eventName string, _ bus.EventHandler) {

	b.lock.Lock()
	defer b.lock.Unlock()

	b.handlersMap[eventName] = slices.FilterNot(b.handlersMap[eventName], func(e bus.EventHandler) bool {
		return false
		// return e == handler TODO solve this comparison
	})
}

func (b *Broker) unsubscribeAll(unsubscribeFuncs []UnsubscribeFunc) {
	for _, unsubscribe := range unsubscribeFuncs {
		unsubscribe()
	}
}

// Handle an event by calling all the handlers subscribed to the event name
// All handlers are called in parallel and the function returns when all handlers have finished
// If any of the handlers returns an error, the function returns an error with all the errors messages
func (b *Broker) Handle(ctx context.Context, ev event.Event) error {

	_, errs := syncs.WaitAll(ctx, slices.Map(b.getHandlers(ev), func(handle bus.EventHandler) syncs.Waitable[any] {
		return func(ctx context.Context) (any, error) {
			err := handle(ctx, ev)

			if err != nil {
				b.logger.Warn(
					"error processing event",
					zap.String("type", ev.Name()),
					zap.Any("id", ev.ID()),
					zap.Error(err),
				)
			}

			return nil, err
		}
	})...)

	filteredErrs := slices.Filter(errs, func(e error) bool { return e != nil })

	if len(filteredErrs) > 0 {
		// TODO check if it can be retried
		// TODO improve semantic

		return fmt.Errorf("errors while procesing the event %s: %v", ev.Name(), filteredErrs)
	}
	return nil
}

func (b *Broker) getHandlers(ev event.Event) []bus.EventHandler {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.handlersMap.EventHandlers(ev.Name())
}

func slicesToSet[T comparable](keys []T) map[T]bool {
	m := make(map[T]bool)
	for _, k := range keys {
		m[k] = true
	}
	return m
}
