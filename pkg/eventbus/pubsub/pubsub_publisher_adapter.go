package pubsub

import (
	"errors"
	"github.com/AltScore/gothic/pkg/eventbus"
	"github.com/modernice/goes/event"
)

// PublisherAdapter is an adapter that allows to use a pubsub.Publisher as an eventbus.Publisher
type PublisherAdapter struct {
	pubSubPublisher *Publisher
}

func NewPublisherAdapter(pubSubPublisher *Publisher) *PublisherAdapter {
	return &PublisherAdapter{pubSubPublisher: pubSubPublisher}
}

func (p *PublisherAdapter) Publish(event1 eventbus.Event, options ...eventbus.Option) error {
	ev, ok := event1.(event.Event)

	if !ok {
		return errors.New("event is not of type *event.Event")
	}

	return p.pubSubPublisher.Publish(ev, options...)
}
