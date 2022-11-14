package main

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/leonsteinhaeuser/observer/v2"
)

// Event with some payload
type Event struct {
	ID      int
	Payload int32
}

// Handler models some long-running component that needs to react to external events
type Handler struct {
	N int32
}

// onEvent processes Event
func (h *Handler) OnEvent(ctx context.Context, message Event) error {
	log.Printf("received message(%d) with payload %d\n", message.ID, message.Payload)
	// do something with event's payload
	atomic.AddInt32(&h.N, message.Payload)
	return nil
}

var obsrv = new(observer.Observer[Event])

func main() {
	h := &Handler{}
	g, ctx := errgroup.WithContext(context.TODO())
	g.Go(observer.Handle[Event](ctx, obsrv, h.OnEvent))
	for id, payload := range []int32{1, 2, 3, 5, -1, 10} {
		obsrv.NotifyAll(Event{id, payload})
	}
	time.Sleep(1 * time.Second) // Allow events to settle
	if err := obsrv.Close(); err != nil {
		log.Printf("closing observer: %+v", err)
	}
	log.Printf("final result: %d", h.N)
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
