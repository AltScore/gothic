package xeh

import (
	"context"
	"errors"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/eventhandler/projector"
	"github.com/looplab/eventhorizon/uuid"
	"go.uber.org/zap"
)

// EntityHealer is a healer of entities of the read model.
// It wraps a Projector to recreate a read model from the event store in case of a failure eh.ErrIncorrectEntityVersion.
type EntityHealer struct {
	*projector.EventHandler

	logger *zap.Logger

	eventStore eh.EventStore
	repo       eh.ReadWriteRepo
}

func NewEntityHealer(logger *zap.Logger, eventStore eh.EventStore, prj projector.Projector, repo eh.ReadWriteRepo, options ...projector.Option) *EntityHealer {

	return &EntityHealer{
		logger:       logger,
		EventHandler: projector.NewEventHandler(prj, repo, options...),
		eventStore:   eventStore,
		repo:         repo,
	}
}

func (h *EntityHealer) HandleEvent(ctx context.Context, event eh.Event) error {
	err := h.EventHandler.HandleEvent(ctx, event)

	if err == nil {
		// No error, return
		return nil
	}

	if !errors.Is(err, eh.ErrIncorrectEntityVersion) {
		// Not an error we can handle, return it
		return err
	}

	if healingErr := h.healEntity(ctx, event); healingErr != nil {
		// Failed to heal the entity, return the original error
		return err
	}

	// Retry the projection
	return h.EventHandler.HandleEvent(ctx, event)
}

func (h *EntityHealer) healEntity(ctx context.Context, event eh.Event) error {
	type_ := event.AggregateType()
	aggId := event.AggregateID()

	h.logger.Warn("Regenerating read models", zap.String("agg_type", type_.String()), zap.String("agg_id", aggId.String()))

	if err := h.removeReadModels(ctx, aggId); err != nil {
		return err
	}

	return h.replayEvents(ctx, type_, aggId)
}

// removeReadModels removes read models for the given aggregate id
func (h *EntityHealer) removeReadModels(ctx context.Context, id uuid.UUID) error {
	err := h.repo.Remove(ctx, id)
	if err == nil || IsEHNotFound(err) {
		return nil
	}
	return err
}

// replayEvents replays events for the given aggregate id from the first event
// Events are not published outside the own event bus
func (h *EntityHealer) replayEvents(ctx context.Context, aggregateType eh.AggregateType, id uuid.UUID) error {
	events, err := h.eventStore.Load(ctx, id)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}

	if events[0].AggregateType() != aggregateType {
		return eh.ErrMismatchedEventAggregateTypes
	}

	eventHandler := h.EventHandler

	for _, event := range events {
		if err := eventHandler.HandleEvent(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
