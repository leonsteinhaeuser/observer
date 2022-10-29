package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/leonsteinhaeuser/observer/v2"
)

const (
	registerClientCount = 3
)

type Event struct {
	ID      int
	Message string
}

var obsrv *observer.Observer[Event] = new(observer.Observer[Event])

func main() {
	wg := sync.WaitGroup{}
	wg.Add(registerClientCount)
	cancelFuns := []observer.CancelFunc{}
	for i := 0; i < registerClientCount; i++ {
		ch, cancelFunc := obsrv.Subscribe()
		cancelFuns = append(cancelFuns, cancelFunc)
		go func(i int, ch <-chan Event) {
			for {
				message := <-ch
				fmt.Printf("Runner: %d\tMessageID: %d\tMessage: %s\n", i, message.ID, message.Message)
				if message.ID == 1 {
					wg.Done()
				}
			}
		}(i, ch)
	}

	fmt.Println("Registered clients:", obsrv.Clients())

	obsrv.NotifyAll(Event{
		ID:      1,
		Message: "Hello World",
	})

	wg.Wait()

	// remove the first client
	err := cancelFuns[0]()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// let's try to remove the first client again
	err = cancelFuns[0]()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// let's send another message
	obsrv.NotifyAll(Event{
		ID:      2,
		Message: "Hello World 2",
	})

	// list clients
	fmt.Println("Registered clients:", obsrv.Clients())

	// deregister all clients
	for _, cancelFunc := range cancelFuns {
		err = cancelFunc()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	// list clients
	fmt.Println("Registered clients:", obsrv.Clients())

	if obsrv.Clients() == 0 {
		fmt.Println("No clients left. ")
		time.Sleep(time.Second * 5)
		return
	}

	select {}
}
