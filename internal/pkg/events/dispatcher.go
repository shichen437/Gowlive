package events

import (
	"container/list"
	"context"
	"sync"

	"github.com/shichen437/gowlive/internal/pkg/interfaces"
)

type dispatcher struct {
	sync.RWMutex
	saver map[EventType]*list.List
}

type Dispatcher interface {
	interfaces.Module
	AddEventListener(eventType EventType, listener *EventListener)
	RemoveEventListener(eventType EventType, listener *EventListener)
	RemoveAllEventListener(eventType EventType)
	DispatchEvent(event *Event)
}

func NewDispatcher(ctx context.Context) Dispatcher {
	return &dispatcher{
		saver: make(map[EventType]*list.List),
	}
}

func (e *dispatcher) DispatchEvent(event *Event) {
	if event == nil {
		return
	}
	e.RLock()
	listeners, ok := e.saver[event.Type]
	if !ok || listeners == nil {
		e.RUnlock()
		return
	}
	hs := make([]*EventListener, 0)
	for e := listeners.Front(); e != nil; e = e.Next() {
		hs = append(hs, e.Value.(*EventListener))
	}
	e.RUnlock()
	go func() {
		for _, h := range hs {
			h.Handler(event)
		}
	}()
}

func (e *dispatcher) AddEventListener(eventType EventType, listener *EventListener) {
	e.Lock()
	defer e.Unlock()
	listeners, ok := e.saver[eventType]
	if !ok || listener == nil {
		listeners = list.New()
		e.saver[eventType] = listeners
	}
	listeners.PushBack(listener)
}

func (e *dispatcher) RemoveEventListener(eventType EventType, listener *EventListener) {
	e.Lock()
	defer e.Unlock()
	listeners, ok := e.saver[eventType]
	if !ok || listeners == nil {
		return
	}
	for e := listeners.Front(); e != nil; e = e.Next() {
		if e.Value == listener {
			listeners.Remove(e)
		}
	}
	if listeners.Len() == 0 {
		delete(e.saver, eventType)
	}
}

func (e *dispatcher) RemoveAllEventListener(eventType EventType) {
	e.Lock()
	defer e.Unlock()
	e.saver = make(map[EventType]*list.List)
}

func (e *dispatcher) Start(ctx context.Context) error {
	return nil
}

func (e *dispatcher) Close(ctx context.Context) {

}
