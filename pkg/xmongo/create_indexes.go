package xmongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// CreateIndexes creates one or more indexes in the provided collection reporting the result in the log
func CreateIndexes(logger *zap.Logger, collection *mongo.Collection, models ...mongo.IndexModel) {
	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	for _, model := range models {
		var ctx context.Context
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		indexName, err := collection.Indexes().CreateOne(ctx, model)

		if err != nil {
			fields := []zap.Field{
				zap.String("collection", collection.Name()),
				zap.Error(err),
			}

			if model.Options.Name != nil {
				fields = append(fields, zap.String("index", *model.Options.Name))
			}

			logger.Warn("Failed to create index", fields...)
		} else {
			logger.Info("Created index", zap.String("collection", collection.Name()), zap.String("index", indexName))
		}

		cancel()
	}
}
