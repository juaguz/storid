package dispatcher

import "context"

type EventHandler interface {
	Handle(ctx context.Context, event string, payload []byte)
}

type EventDispatcher interface {
	Register(ctx context.Context, eventName string, handler EventHandler)
	Dispatch(ctx context.Context, event string, payload []byte)
}

type Dispatcher struct {
	handlers map[string][]EventHandler
	async    bool
}

func NewSimpleEventDispatcher(async bool) *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]EventHandler),
		async:    async,
	}
}

func (d *Dispatcher) Register(ctx context.Context, eventName string, handler EventHandler) {
	if _, exists := d.handlers[eventName]; !exists {
		d.handlers[eventName] = []EventHandler{}
	}
	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

func (d *Dispatcher) Dispatch(ctx context.Context, event string, payload []byte) {
	if handlers, exists := d.handlers[event]; exists {
		for _, handler := range handlers {
			if d.async {
				go handler.Handle(ctx, event, payload)
				continue
			}
			handler.Handle(ctx, event, payload)
		}
	}
}
