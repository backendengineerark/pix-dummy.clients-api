package events

import (
	"context"
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
	SetPayload(payload interface{})
}

type EventHandlerInterface interface {
	Handle(ctx context.Context, event EventInterface, wg *sync.WaitGroup)
}

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(ctx context.Context, event EventInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
}
