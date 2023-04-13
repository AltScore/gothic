package xmongo

import (
	"context"
	"errors"
	"github.com/AltScore/gothic/pkg/entity"
	"github.com/AltScore/gothic/pkg/xerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	dbopts "go.mongodb.org/mongo-driver/mongo/options"
)

var optimisticLockingError = errors.New("optimistic locking failed")

type MongoErrorConverterFunc func(err error, entity string, keyFmt string, args ...interface{}) error

type Option[Entity entity.Entity, Dto any] func(*Repository[Entity, Dto])

type Repository[Entity entity.Entity, Dto any] struct {
	c                 *mongo.Collection
	convertMongoError MongoErrorConverterFunc
	entityName        string
}

func NewRepository[Entity entity.Entity, Dto any](client *mongo.Client, databaseName string, collectionName string, options ...Option[Entity, Dto]) Repository[Entity, Dto] {
	r := Repository[Entity, Dto]{
		c:                 client.Database(databaseName).Collection(collectionName),
		convertMongoError: ConvertMongoError,
		entityName:        collectionName,
	}
	for _, option := range options {
		option(&r)
	}
	return r
}

func (r *Repository[Entity, Dto]) FindOne(ctx context.Context, filter interface{}) (Entity, error) {
	var e Entity
	if err := r.c.FindOne(ctx, filter).Decode(&e); err != nil {
		return e, r.convertMongoError(err, r.entityName, "Filter %s", filter)
	}
	return e, nil
}

func (r *Repository[Entity, Dto]) Find(ctx context.Context, filter interface{}) ([]Entity, error) {
	cursor, err := r.c.Find(ctx, filter)
	if err != nil {
		return nil, r.convertMongoError(err, r.entityName, "Filter %s", filter)
	}
	var entities []Entity
	if err := cursor.All(ctx, &entities); err != nil {
		return entities, r.convertMongoError(err, r.entityName, "Filter %s", filter)
	}

	return entities, nil
}

// FindPage

func (r *Repository[Entity, Dto]) FindById(ctx context.Context, id interface{}) (Entity, error) {
	return r.FindOne(ctx, bson.M{"_id": id})
}

func (r *Repository[Entity, Dto]) FindAll(ctx context.Context) ([]Entity, error) {
	return r.Find(ctx, bson.M{})
}

func (r *Repository[Entity, Dto]) Count(ctx context.Context) (int64, error) {
	return r.c.CountDocuments(ctx, bson.M{})
}

func (r *Repository[Entity, Dto]) FindByIds(ctx context.Context, ids ...interface{}) ([]Entity, error) {
	return r.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
}

func (r *Repository[Entity, Dto]) Store(ctx context.Context, entity Entity) error {
	return nil
}

type Updater[Entity any] func(ctx context.Context, entity Entity) (Entity, error)

// Update updates an entity with the provided updater function
// The updater function is passed the entity to update and should return the updated entity
// The update function will retry 5 times if an optimistic locking error occurs
// Sample use
//
//	func (e *SomeCommandHandler) SetName(ctx context.Context, cmd SetNameCommand) {
//	    e.repository.Update(ctx, e.Id(), func(ctx context.Context, e Entity) (Entity, error) {
//	        err := e.ChangeName(cmd) // this changes the name and creates NameChangedEvent
//	        return e, err
//	    })
//	}
//
// Repository will update the entity and publish the NameChangedEvent
func (r *Repository[Entity, Dto]) Update(ctx context.Context, id interface{}, updater Updater[Entity]) error {
	for i := 0; i < 5; i++ {
		err := r.updateOnce(ctx, id, updater)
		if err != optimisticLockingError {
			return err
		}
		// Should it be delayed?
	}

	return xerrors.NewOptimisticLockingError(r.entityName, "too many retries to update", "id %v", id)
}

func (r *Repository[Entity, Dto]) updateOnce(ctx context.Context, id interface{}, updater Updater[Entity]) error {
	e, err := r.FindById(ctx, id)
	if err != nil {
		return err
	}

	updated, err := updater(ctx, e)

	if err != nil {
		return err
	}

	f := bson.M{"_id": id, "version": e.Version()}

	result, err := r.c.UpdateOne(ctx, f, updated)

	if err != nil {
		return r.convertMongoError(err, r.entityName, "Id %s", id)
	}

	if result.MatchedCount == 0 {
		return optimisticLockingError
	}

	return nil
}

// WithRegistry provides a registry to encode / decode entities
func WithRegistry[Entity any, Dto any](registry *bsoncodec.Registry) func(*Repository[Entity, Dto]) {
	return func(r *Repository[Entity, Dto]) {
		r.c = r.c.Database().Collection(r.c.Name(), dbopts.Collection().SetRegistry(registry))
	}
}

// WithErrorConverter provides a custom error converter
func WithErrorConverter[Entity any, Dto any](converter MongoErrorConverterFunc) func(*Repository[Entity, Dto]) {
	return func(r *Repository[Entity, Dto]) {
		r.convertMongoError = converter
	}
}
