package observer

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-multierror"
)

var ErrClientAlreadyDeRegistered = errors.New("client already de-registered")

type CancelFunc func() error

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
	mu            sync.RWMutex
	n             int64
	numDeleted    int64
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
	return listenCh, func() error {
		return o.deleteClient(n)
	}
}

func (o *Observer[T]) deleteClient(key any) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	c, ok := o.clients.LoadAndDelete(key)
	if !ok {
		return ErrClientAlreadyDeRegistered
	}
	atomic.AddInt64(&o.numDeleted, 1)
	close(c.(chan T))
	return nil
}

// NotifyAll notifies all registered clients.
func (o *Observer[T]) NotifyAll(data T) {
	o.clients.Range(func(key, _ any) bool {
		go func(key any) {
			o.mu.RLock()
			defer o.mu.RUnlock()
			client, ok := o.clients.Load(key)
			if !ok {
				return
			}
			select {
			case client.(chan T) <- data:
				// the message was sent successfully
				return
			case <-time.After(o.notifyTimeout()):
				// client is not responding
				return
			}
		}(key)
		return true
	})
}

// Clients returns the number of registered clients.
func (o *Observer[t]) Clients() int64 {
	return o.n - o.numDeleted
}

// Handle builds repetitive message consumer using provided handler function h.
// Returned func() error value is suitable to run in errrgroup's Go() method.
func Handle[T any](ctx context.Context, o Observable[T], h func(context.Context, T) error) func() error {
	msgs, unsub := o.Subscribe()
	return func() error {
		defer unsub()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case upd, ok := <-msgs:
				if !ok {
					return nil
				}
				if err := h(ctx, upd); err != nil {
					return err
				}
			}
		}
	}
}

// Close disconnects all clients from the observer
func (o *Observer[T]) Close() error {
	var res error
	o.clients.Range(func(key, _ any) bool {
		if err := o.deleteClient(key); err != nil {
			res = multierror.Append(res, err)
		}
		return true
	})
	return res
}
