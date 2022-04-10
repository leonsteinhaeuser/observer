# observer

This library implements the pub/sub pattern in a generic way. It uses Go's generic types to declare the type of the client identifier as well as the type of the event.

## Usage

```bash
go get github.com/leonsteinhaeuser/observer
```

## Example

```go
package main

import (
    "fmt"
    "github.com/leonsteinhaeuser/observer"
)

type Event struct {
    ID int
    Message string
}

var (
    obsrv observer.Observable[int, Event] = observer.NewObserver[int, Event]()
)

func main() {
    obsrv.RegisterClient(1, make(chan Event))
    obsrv.NotifyAll(Event{
        ID:      i,
        Message: fmt.Sprintf("Message with ID %d", i),
    })
    fmt.Println("Registered Clients: ", obsrv.Clients())
    err := obsrv.DeRegisterClient(1)
    if err != nil {
        fmt.Println(err)
    }
}
```

If you would like to see a more detailed example, please take a look at the [observer](_example/main.go) example.
