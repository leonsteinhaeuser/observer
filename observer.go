package observer

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrClientAlreadyDeRegistered = errors.New("client already de-registered")

type CancelFunc func()

// Observable
type Observable[T any] interface {
	// Subscribe registers a client to the observable.
	Subscribe() (<-chan T, CancelFunc)
	// NotifyAll notifies all registered clients.
	NotifyAll(data T)
}

// Observer offers the possibility to notify all registered clients.
// Since the client map must be initialized, it is not possible to use this structure directly.
// Use NewObserver instead.
type Observer[T any] struct {
	NotifyTimeout time.Duration
	clients       sync.Map
	n             int64
}

func (o *Observer[T]) notifyTimeout() time.Duration {
	if o.NotifyTimeout != 0 {
		return o.NotifyTimeout
	}
	return time.Second * 5
}

// Subscribe registers a client to the observer and returns a channel to receive notifications.
// The returned CancelFunc can be used to de-register the client.
func (o *Observer[T]) Subscribe() (<-chan T, CancelFunc) {
	n := atomic.AddInt64(&o.n, 1)
	listenCh := make(chan T)
	o.clients.Store(n, listenCh)
	return listenCh, func() {
		o.deleteClient(n)
	}
}

func (o *Observer[T]) deleteClient(uid int64) {
	c, ok := o.clients.LoadAndDelete(uid)
	if !ok {
		return
	}
	close(c.(chan T))
}

// NotifyAll notifies all registered clients.
func (o *Observer[T]) NotifyAll(data T) {
	o.clients.Range(func(_, value any) bool {
		go func(client chan T) {
			select {
			case client <- data:
				// the message was sent successfully
				return
			case <-time.After(o.notifyTimeout()):
				// client is not responding
				return
			}
		}(value.(chan T))
		return true
	})
}
