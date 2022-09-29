package codec

import (
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

type Unmarshaller func(data []byte, value reflect.Value) error

type Registry interface {
	ForName(name string) (Unmarshaller, bool)
}

type registry struct {
	mu        sync.RWMutex
	factories map[string]Unmarshaller
}

var defaultRegistry = New()

func New() *registry {
	return &registry{
		factories: make(map[string]Unmarshaller),
	}
}

func Default() *registry {
	return defaultRegistry
}

func ForName(name string) (Unmarshaller, bool) {
	return Default().ForName(name)
}

func (r *registry) ForName(name string) (Unmarshaller, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.factories[name]
	return f, ok
}

func (r *registry) Register(name string, f Unmarshaller) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[name] = f
}

func Register[T any](name string) {

	Default().Register(name, func(data []byte, value reflect.Value) error {
		t := new(T)

		if err := bson.Unmarshal(data, &t); err != nil {
			return err
		}

		newValue := reflect.ValueOf(t)
		elem := value.Elem()
		elem.Set(newValue)

		return nil
	})
	//		return reflect.New(tt).Interface()
}
