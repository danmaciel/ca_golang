package events

import (
	"errors"
	"sync"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

// EventDispatcher protects its subscription registry with a RWMutex.
// Dispatch captures a handler snapshot and releases the lock before invoking
// handlers; subscriptions changed during a dispatch affect later dispatches.
type EventDispatcher struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ev *EventDispatcher) Dispatch(event EventInterface) error {
	ev.mu.RLock()
	handlers := append([]EventHandlerInterface(nil), ev.handlers[event.GetName()]...)
	ev.mu.RUnlock()

	var wg sync.WaitGroup
	for _, handler := range handlers {
		wg.Add(1)
		go handler.Handle(event, &wg)
	}
	wg.Wait()
	return nil
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	ed.mu.Lock()
	defer ed.mu.Unlock()
	for _, h := range ed.handlers[eventName] {
		if h == handler {
			return ErrHandlerAlreadyRegistered
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	ed.mu.RLock()
	defer ed.mu.RUnlock()
	for _, h := range ed.handlers[eventName] {
		if h == handler {
			return true
		}
	}
	return false
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	ed.mu.Lock()
	defer ed.mu.Unlock()
	for i, h := range ed.handlers[eventName] {
		if h == handler {
			ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
			return nil
		}
	}
	return nil
}

func (ed *EventDispatcher) Clear() {
	ed.mu.Lock()
	defer ed.mu.Unlock()
	ed.handlers = make(map[string][]EventHandlerInterface)
}
