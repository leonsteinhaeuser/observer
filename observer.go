package observer

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrClientNotFound = errors.New("client not found")
)

// Observable
type Observable[i comparable, t any] interface {
	// RegisterClient registers a client by its id.
	RegisterClient(id i, client chan t)
	// DeRegisterClient deregisters a client by its id.
	DeRegisterClient(id i) error
	// NotifyAll notifies all registered clients.
	NotifyAll(data t)
	// NotifyClient notifies a single registered client by its id.
	NotifyClient(id i, data t) error
	// Clients returns the number of registered clients.
	Clients() int64
}

// Observer offers the possibility to notify all registered clients.
// Since the client map must be initialized, it is not possible to use this structure directly.
// Use NewObserver instead.
type Observer[i comparable, t any] struct {
	lock    sync.Mutex
	clients map[i]chan t
}

// NewObserver initializes a new Observer.
func NewObserver[i comparable, t any]() *Observer[i, t] {
	return &Observer[i, t]{
		clients: make(map[i]chan t),
	}
}

// RegisterClient registers a client by its id.
func (o *Observer[i, t]) RegisterClient(id i, client chan t) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.clients[id] = client
}

// DeRegisterClient deregisters a client by its id.
// If the client is not registered, an error is returned.
// Note that the channel registered for the client will not be closed.
func (o *Observer[i, t]) DeRegisterClient(id i) error {
	o.lock.Lock()
	defer o.lock.Unlock()
	if _, ok := o.clients[id]; !ok {
		return fmt.Errorf("%w: %v", ErrClientNotFound, id)
	}
	delete(o.clients, id)
	return nil
}

// NotifyAll notifies all registered clients.
func (o *Observer[i, t]) NotifyAll(data t) {
	for _, client := range o.clients {
		client <- data
	}
}

// NotifyClient notifies a single registered client by its id.
func (o *Observer[i, t]) NotifyClient(id i, data t) error {
	if _, ok := o.clients[id]; !ok {
		return fmt.Errorf("%w: %v", ErrClientNotFound, id)
	}
	o.clients[id] <- data
	return nil
}

// Clients returns the number of registered clients.
func (o *Observer[i, t]) Clients() int64 {
	return int64(len(o.clients))
}
