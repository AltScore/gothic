package xbson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_mongo_bson_honor_json_tags(t *testing.T) {
	saved := bson.DefaultRegistry
	defer func() { bson.DefaultRegistry = saved }()

	NewBsonRegistryBuilder().Build()

	type Test struct {
		FirstName string `json:"first_name"`
		Age       int    `json:"yearsPast"`
	}

	test := Test{FirstName: "John", Age: 42}

	bsonTest, err := bson.Marshal(&test)

	assert.Nil(t, err)

	var unmarshalled map[string]interface{}

	err = bson.Unmarshal(bsonTest, &unmarshalled)

	t.Log(unmarshalled)

	assert.Nil(t, err)
	assert.Equal(t, "John", unmarshalled["first_name"])
	assert.Equal(t, int32(42), unmarshalled["yearsPast"])
}
