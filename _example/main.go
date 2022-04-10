package main

import (
	"fmt"
	"time"

	"github.com/leonsteinhaeuser/observer"
)

type Event struct {
	ID      int
	Message string
}

var (
	obsrv observer.Observable[int, Event] = observer.NewObserver[int, Event]()
)

func main() {
	// mock clients and their channels
	clients := make(map[int]chan Event)
	for i := 0; i < 5; i++ {
		clntChan := make(chan Event)
		go func(id int) {
			// example work
			for {
				select {
				case evt := <-clntChan:
					fmt.Printf("Client %d received: %+v\n", id, evt)
				}
			}
		}(i)
		clients[i] = clntChan
	}

	// register clients
	for id, client := range clients {
		obsrv.RegisterClient(id, client)
	}

	// send a bunch of events
	go func() {
		for i := 0; i < 10; i++ {
			obsrv.NotifyAll(Event{
				ID:      i,
				Message: fmt.Sprintf("Message with ID %d", i),
			})
		}
	}()

	// private function to deregister a client
	fncDeregister := func(numbers int) {
		counter := 0
		for key, _ := range clients {
			if counter == numbers {
				break
			}
			counter++

			err := obsrv.DeRegisterClient(key)
			if err != nil {
				fmt.Println(err)
			}
			delete(clients, key)
		}
	}

	fmt.Println("========== REGISTERED CLIENTS ==========: ", obsrv.Clients())

	time.Sleep(time.Second * 5)

	fmt.Println("========== REGISTERED CLIENTS ==========: ", obsrv.Clients())
	// remove two clients
	fncDeregister(2)

	fmt.Println("========== REGISTERED CLIENTS ==========: ", obsrv.Clients())

	time.Sleep(time.Second * 5)

	// remove additional 5 clients
	fncDeregister(5)

	// all clients have been removed
	fmt.Println("========== REGISTERED CLIENTS ==========: ", obsrv.Clients())
}
