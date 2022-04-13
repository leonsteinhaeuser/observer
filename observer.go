package observer

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrClientAlreadyDeRegistered = errors.New("client already de-registered")
)

type CancelFunc func() error

// Observable
type Observable[t any] interface {
	// Subscribe registers a client to the observable.
	Subscribe() (<-chan t, CancelFunc)
	// NotifyAll notifies all registered clients.
	NotifyAll(data t)
	// Clients returns the number of registered clients.
	Clients() int64
}

// Observer offers the possibility to notify all registered clients.
// Since the client map must be initialized, it is not possible to use this structure directly.
// Use NewObserver instead.
type Observer[t any] struct {
	lock    sync.Mutex
	clients map[string]chan t
}

// NewObserver initializes a new Observer.
func NewObserver[t any]() *Observer[t] {
	return &Observer[t]{
		clients: make(map[string]chan t),
	}
}

// Subscribe registers a client to the observer and returns a channel to receive notifications.
// The returned CancelFunc can be used to de-register the client.
func (o *Observer[t]) Subscribe() (<-chan t, CancelFunc) {
	o.lock.Lock()
	defer o.lock.Unlock()
	uid := uuid.NewString()
	listenCh := make(chan t)
	o.clients[uid] = listenCh
	return listenCh, func() error {
		return o.deleteClient(uid)
	}
}

func (o *Observer[t]) deleteClient(uid string) error {
	if _, ok := o.clients[uid]; !ok {
		return ErrClientAlreadyDeRegistered
	}
	close(o.clients[uid])
	delete(o.clients, uid)
	return nil
}

// NotifyAll notifies all registered clients.
func (o *Observer[t]) NotifyAll(data t) {
	for _, client := range o.clients {
		go func(client chan t) {
			select {
			case client <- data:
				// the message was sent successfully
				return
			case <-time.After(time.Second * 5):
				// client is not responding
				return
			}
		}(client)
	}
}

// Clients returns the number of registered clients.
func (o *Observer[t]) Clients() int64 {
	return int64(len(o.clients))
}
