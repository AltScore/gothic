package xmongo

import (
	"github.com/AltScore/gothic/pkg/xpaging"
	"go.mongodb.org/mongo-driver/bson"
)

func ConvertSortOptionsToMongo(defaultField string, defaultDirection xpaging.SortDirection, options xpaging.SortOptions) bson.D {
	var sort bson.D
	if len(options) > 0 {

		for _, option := range options {
			sort = append(sort, bson.E{Key: option.FieldName, Value: int(option.Direction)})
		}
	} else {
		sort = append(sort, bson.E{Key: defaultField, Value: defaultDirection})
	}

	return sort
}
