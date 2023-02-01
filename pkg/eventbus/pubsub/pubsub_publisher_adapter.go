package pubsub

import (
	"errors"
	"github.com/AltScore/gothic/pkg/es/event"
	"github.com/AltScore/gothic/pkg/eventbus"
)

// PublisherAdapter is an adapter that allows to use a pubsub.Publisher as an eventbus.Publisher
type PublisherAdapter struct {
	pubSubPublisher *Publisher
}

func NewPublisherAdapter(pubSubPublisher *Publisher) *PublisherAdapter {
	return &PublisherAdapter{pubSubPublisher: pubSubPublisher}
}

func (p *PublisherAdapter) Publish(event1 eventbus.Event, options ...eventbus.Option) error {
	ev, ok := event1.(event.IEvent)

	if !ok {
		return errors.New("event is not of type *event.Event")
	}

	return p.pubSubPublisher.Publish(ev, options...)
}
