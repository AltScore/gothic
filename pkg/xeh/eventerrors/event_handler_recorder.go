package eventerrors

import (
	"context"
	"github.com/AltScore/gothic/v2/pkg/ids"
	"github.com/AltScore/gothic/v2/pkg/xcontext"
	"github.com/AltScore/gothic/v2/pkg/xerrors"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/repo/mongodb"
	"github.com/looplab/eventhorizon/uuid"
	"github.com/totemcaf/gollections/ptrs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

// RetriableError is a marker to identify an error that can be retried.
type RetriableError interface {
	IsRetriable() bool
}

type EventError struct {
	Id        ids.Id    `bson:"_id"`
	CreatedAt time.Time `bson:"createdAt"`
	UserId    *ids.Id   `bson:"userId"` // The user that caused the error, if present
	Tenant    *string   `bson:"tenant"` // The tenant that caused the error, if present
	Err       error     `bson:"error"`  // The error returned by the event handler
	Event     eh.Event  `bson:"event"`  // The event that caused the error
	Host      string    `bson:"host"`   // This is the machine name
}

var _ eh.Entity = (*EventError)(nil)

func (e EventError) EntityID() uuid.UUID { return e.Id }

// EventHandlerErrorRecorder is a wrapper to a EventHandler to catch the returned errors from a target event handler,
// and record them to process them.
type EventHandlerErrorRecorder struct {
	logger *zap.Logger
	target eh.EventBus
	store  eh.WriteRepo
}

// NewEventHandlerErrorRecorder returns a new instance of EventHandlerErrorRecorder using the mongo client to
// record the events in the indicated database and collection.
func NewEventHandlerErrorRecorder(logger *zap.Logger, client *mongo.Client, databaseName, collectionName string, target eh.EventBus) (*EventHandlerErrorRecorder, error) {
	xerrors.EnsureNotEmpty(logger, "logger")
	xerrors.EnsureNotEmpty(client, "client")
	xerrors.EnsureNotEmpty(databaseName, "databaseName")
	xerrors.EnsureNotEmpty(collectionName, "collectionName")
	xerrors.EnsureNotEmpty(target, "target")

	logger = logger.Named(target.HandlerType().String() + " recorder")

	logger.Info("creating event handler error recorder", zap.String("databaseName", databaseName), zap.String("collectionName", collectionName))

	store, err := mongodb.NewRepoWithClient(client, databaseName, collectionName)
	if err != nil {
		return nil, err
	}

	return &EventHandlerErrorRecorder{logger: logger, target: target, store: store}, nil
}

var _ eh.EventBus = (*EventHandlerErrorRecorder)(nil)

// AddHandler implements the AddHandler method of the EventBus interface.
func (e *EventHandlerErrorRecorder) AddHandler(ctx context.Context, matcher eh.EventMatcher, handler eh.EventHandler) error {
	return e.target.AddHandler(ctx, matcher, e.wrap(handler))
}

// HandlerType implements the HandlerType method of the EventBus interface.
func (e *EventHandlerErrorRecorder) HandlerType() eh.EventHandlerType {
	return e.target.HandlerType()
}

// HandleEvent implements the HandleEvent method of the EventBus interface.
func (e *EventHandlerErrorRecorder) HandleEvent(ctx context.Context, event eh.Event) error {
	return e.target.HandleEvent(ctx, event)
}

// Errors implements the Errors method of the EventBus interface.
func (e *EventHandlerErrorRecorder) Errors() <-chan error {
	return e.target.Errors()
}

// Close implements the Close method of the EventBus interface.
func (e *EventHandlerErrorRecorder) Close() error {
	return e.target.Close()
}

// wrap wraps the handler to catch the returned errors.
func (e *EventHandlerErrorRecorder) wrap(handler eh.EventHandler) eh.EventHandler {
	return wrapper{handler: handler}
}

// persistError persists the error in the database for later processing.
// Current user and tenant in context are recorded.
func (e *EventHandlerErrorRecorder) persistError(ctx context.Context, event eh.Event, err error) error {
	userId, tenant := e.getUserAndTenant(ctx)

	eventError := EventError{
		Id:        ids.New(),
		CreatedAt: event.Timestamp(),
		UserId:    userId,
		Tenant:    tenant,
		Err:       err,
		Event:     event,
	}

	if err := e.store.Save(ctx, &eventError); err != nil {
		e.logger.Error("could not save event error", zap.Error(err))
		// TODO send error in error channel

		return err
	}

	return nil
}

func (e *EventHandlerErrorRecorder) getUserAndTenant(ctx context.Context) (*ids.Id, *string) {
	if user, err := xcontext.GetUser(ctx); err == nil {
		id := user.Id()
		return &id, ptrs.Ptr(user.Tenant())
	}

	if t, found := xcontext.GetTenant(ctx); found {
		return nil, &t
	}

	return nil, nil
}

// wrapper is a wrapper to a EventHandler to catch the returned errors from the target event handler,
type wrapper struct {
	handler  eh.EventHandler
	recorder *EventHandlerErrorRecorder
}

var _ eh.EventHandler = (*wrapper)(nil)

func (w wrapper) HandlerType() eh.EventHandlerType {
	return w.handler.HandlerType()
}

func (w wrapper) HandleEvent(ctx context.Context, event eh.Event) error {
	err := w.handler.HandleEvent(ctx, event)
	if err == nil {
		return nil
	}
	if retryable, ok := err.(RetriableError); ok && retryable.IsRetriable() {
		// Error is retriable, so we return it to indicate that the event was not handled.
		return err
	}

	if err := w.recorder.persistError(ctx, event, err); err != nil {
		// Failed to persist the error, so we return it to indicate that the event was not handled.
		return err
	}

	// Error is not retriable, so we return nil to indicate that the event was handled.
	return nil
}
