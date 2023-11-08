package xrepo

import (
	"sync"
)

// InMemoryRepo is a simple in-memory repository implementation.
// It is thread-safe, so it should be wrapped in a mutex.
type InMemoryRepo[Entry any] struct {
	entries map[string]Entry
	lock    sync.RWMutex
	keyFn   func(Entry) string
}

func NewInMemoryRepo[Entry any](keyFn func(Entry) string) *InMemoryRepo[Entry] {
	return &InMemoryRepo[Entry]{
		entries: make(map[string]Entry),
		keyFn:   keyFn,
	}
}

func (r *InMemoryRepo[Entry]) FindByKey(key string) (Entry, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	entry, ok := r.entries[key]
	return entry, ok
}

func (r *InMemoryRepo[Entry]) Update(key string, updater func(entry Entry, found bool) (Entry, error)) (Entry, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	entry, found := r.entries[key]

	newEntry, err := updater(entry, found)

	if err == nil {
		if r.keyFn(newEntry) == "" {
			r.internalDelete(key)
		} else {
			r.internalPut(newEntry)
		}
	}

	return newEntry, err
}

func (r *InMemoryRepo[Entry]) Store(entry Entry) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.internalPut(entry)
}

func (r *InMemoryRepo[Entry]) internalPut(entry Entry) {
	r.entries[r.keyFn(entry)] = entry
}

func (r *InMemoryRepo[Entry]) Delete(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.internalDelete(key)
}

func (r *InMemoryRepo[Entry]) internalDelete(key string) {
	delete(r.entries, key)
}

func (r *InMemoryRepo[Entry]) Find(filter func(Entry) bool) []Entry {
	r.lock.RLock()
	defer r.lock.RUnlock()

	entries := make([]Entry, 0)
	for _, entry := range r.entries {
		if filter(entry) {
			entries = append(entries, entry)
		}
	}

	return entries
}
