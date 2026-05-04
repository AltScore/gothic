package xbson

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	registrars []Registrar

	// DefaultRegistry holds the registry produced by BuildDefaultRegistry. In
	// mongo-driver v2 there is no global default registry on the bson package;
	// consumers must pass this registry explicitly via
	// options.Client().SetRegistry(...) or options.Database().SetRegistry(...).
	DefaultRegistry *bson.Registry
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
// decoders, plus all registered registrars.
func BuildRegistry() *bson.Registry {
	registry := bson.NewRegistry()

	for _, registrar := range registrars {
		registrar.Register(registry)
	}

	return registry
}

// BuildDefaultRegistry builds the registry and stores it in xbson.DefaultRegistry.
// In mongo-driver v2 there is no settable global default registry, so consumers
// must pass xbson.DefaultRegistry explicitly to their mongo client/database.
func BuildDefaultRegistry() {
	DefaultRegistry = BuildRegistry()
}
