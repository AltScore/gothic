package xbson

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v3"
)

func Dump(t *testing.T, title string, bytes []byte) {
	var bsonMap bson.M

	err := bson.Unmarshal(bytes, &bsonMap)
	require.NoError(t, err)

	// Convert object to YAML
	yamlBytes, err := yaml.Marshal(bsonMap)

	require.NoError(t, err)

	yamlString := string(yamlBytes)
	fmt.Println(title)
	fmt.Println(yamlString)

	fmt.Println(bsonMap)
}
