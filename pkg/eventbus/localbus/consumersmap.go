package localbus

import (
	"github.com/AltScore/gothic/pkg/eventbus"
	"sync"
)

type consumersMap struct {
	consumers map[eventbus.EventName][]eventbus.EventConsumer
	lock      sync.RWMutex
}

func (l *consumersMap) addConsumer(eventName eventbus.EventName, consumer eventbus.EventConsumer) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.consumers == nil {
		l.consumers = make(map[eventbus.EventName][]eventbus.EventConsumer)
	}

	l.consumers[eventName] = append(l.consumers[eventName], consumer)
}

func (l *consumersMap) getConsumers(eventName eventbus.EventName) []eventbus.EventConsumer {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.consumers[eventName]
}
