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

func main() {
	observer := observer.NewObserver[string, Event]()

	clients := make(map[string]chan Event)
	clients = map[string]chan Event{
		"client1": make(chan Event),
		"client2": make(chan Event),
		"client3": make(chan Event),
		"client4": make(chan Event),
		"client5": make(chan Event),
	}

	for key, client := range clients {
		go func(key string, client chan Event) {
			fmt.Println("Registering client:", key)
			observer.RegisterClient(key, client)

			for {
				data := <-client
				fmt.Println("Client:", key, "received:", data)
			}
		}(key, client)
	}

	go func() {
		for i := 0; i <= 2; i++ {
			go func(idx int) {
				for j := 0; j <= 5; j++ {
					time.Sleep(time.Second * 2)
					observer.NotifyAll(Event{
						ID:      idx,
						Message: fmt.Sprintf("custom: %d=%d", idx, j),
					})
				}
			}(i)
		}
	}()

	fmt.Println("========== 1 REGISTERED CLIENTS ==========: ", observer.Clients())

	time.Sleep(time.Second * 5)

	fmt.Println("========== 2 REGISTERED CLIENTS ==========: ", observer.Clients())
	// remove two clients
	counter := 0
	for key, _ := range clients {
		if counter == 2 {
			break
		}
		counter++

		err := observer.DeRegisterClient(key)
		if err != nil {
			fmt.Println(err)
		}
		delete(clients, key)
	}
	fmt.Println("========== 3 REGISTERED CLIENTS ==========: ", observer.Clients())

	time.Sleep(time.Second * 5)

	// remove additional 5 clients
	counter = 0
	for key, _ := range clients {
		if counter == 5 {
			break
		}
		counter++

		err := observer.DeRegisterClient(key)
		if err != nil {
			fmt.Println(err)
		}
		delete(clients, key)
	}

	// all clients have been removed
	fmt.Println("========== 4 REGISTERED CLIENTS ==========: ", observer.Clients())
}
