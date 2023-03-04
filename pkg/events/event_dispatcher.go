package events

import (
	"context"
	"errors"
	"sync"
)

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return errors.New("handler already registered")
			}
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Dispatch(ctx context.Context, event EventInterface) error {
	if handlers, ok := ed.handlers[event.GetName()]; ok {
		wg := sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(ctx, event, &wg)
		}
	}
	return nil
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}
