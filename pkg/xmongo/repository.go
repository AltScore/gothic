package xmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	dbopts "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoErrorConverterFunc func(err error, entity string, keyFmt string, args ...interface{}) error

type Option[Entity any, Dto any] func(*Repository[Entity, Dto])

type Repository[Entity any, Dto any] struct {
	c                 *mongo.Collection
	convertMongoError MongoErrorConverterFunc
	entityName        string
}

func NewRepository[Entity, Dto any](client *mongo.Client, databaseName string, collectionName string, options ...Option[Entity, Dto]) Repository[Entity, Dto] {
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
	var entity Entity
	if err := r.c.FindOne(ctx, filter).Decode(&entity); err != nil {
		return entity, r.convertMongoError(err, r.entityName, "Filter %s", filter)
	}
	return entity, nil
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
