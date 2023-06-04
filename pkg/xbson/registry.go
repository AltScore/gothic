package xbson

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

var (
	registrars []Registrar
)

// Register registers a Registrar to the list of registrars.
// If the Registrar is already registered, it will not be registered again.
// This function is not thread-safe, and it is typically called from init() functions.
func Register(registrar Registrar) {
	if !IsAlreadyRegistered(registrar) {
		registrars = append(registrars, registrar)
	}

}

// IsAlreadyRegistered checks if a Registrar is already registered.
func IsAlreadyRegistered(registrar Registrar) bool {
	for _, r := range registrars {
		if r == registrar {
			return true
		}
	}

	return false
}

// BuildRegistry creates a new registry configured with the default encoders and
// decoders from the bsoncodec.DefaultValueEncoders and bsoncodec.DefaultValueDecoders types, the
// PrimitiveCodecs type in this package, and all registered registrars.
func BuildRegistry() *bsoncodec.Registry {
	builder := bson.NewRegistryBuilder()

	for _, registrar := range registrars {
		registrar.Register(builder)
	}

	return builder.Build()
}

// BuildDefaultRegistry builds the default registry to be used by the mongo driver
// Previous registries are discarded
func BuildDefaultRegistry() {
	bson.DefaultRegistry = BuildRegistry()
}
