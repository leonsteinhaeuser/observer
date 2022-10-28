package main

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/leonsteinhaeuser/observer/v2"
	"golang.org/x/sync/errgroup"
)

// Event with some payload
type Event struct {
	ID      int
	Payload int32
}

// Handler models some long-running component that needs to react to external events
type Handler struct {
	N   int32
	g   *errgroup.Group
	cfs []observer.CancelFunc
	wg  sync.WaitGroup
}

// onEvent processes Event
func (h *Handler) onEvent(ctx context.Context, message Event) error {
	log.Printf("received message(%d) with payload %d\n", message.ID, message.Payload)
	// do something with event's payload
	atomic.AddInt32(&h.N, message.Payload)
	return nil
}

// Connect Handler to the external event source (Observable)
func (h *Handler) Connect(o observer.Observable[Event]) error {
	g, ctx := errgroup.WithContext(context.Background())
	h.g = g
	f, cf := observer.Handle(ctx, o, h.onEvent)
	g.Go(f)
	h.cfs = append(h.cfs, cf)
	return nil
}

// Close disconnects external events and may also clean up other dependencies
// after the call to h.g.Wait() returns.
func (h *Handler) Close() error {
	var res error
	for _, cf := range h.cfs {
		if err := cf(); err != nil {
			res = multierror.Append(res, err)
		}
	}
	if err := h.g.Wait(); err != nil {
		res = multierror.Append(res, err)
	}
	// it's now safe to disconnect database etc since event handlers are
	// terminated
	return res
}

var obsrv = new(observer.Observer[Event])

func main() {
	h := &Handler{}
	if err := h.Connect(obsrv); err != nil {
		log.Fatal(err)
	}
	for id, payload := range []int32{1, 2, 3, 5, -1, 10} {
		obsrv.NotifyAll(Event{id, payload})
	}
	// allow events to settle
	time.Sleep(1 * time.Second)
	log.Printf("final result: %d", h.N)
	if err := h.Close(); err != nil {
		log.Fatal(err)
	}
}
