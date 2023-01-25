package localbus

import (
	"github.com/AltScore/gothic/pkg/eventbus"
	"sync"
)

type handlersMap struct {
	handlers map[string][]eventbus.EventHandler
	lock     sync.RWMutex
}

func (l *handlersMap) addHandler(eventName string, listener eventbus.EventHandler) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.handlers == nil {
		l.handlers = make(map[string][]eventbus.EventHandler)
	}

	l.handlers[eventName] = append(l.handlers[eventName], listener)
}

func (l *handlersMap) getHandlers(eventName string) []eventbus.EventHandler {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.handlers[eventName]
}
