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
    ID      int
    Message string
}

var (
    obsrv observer.Observable[Event] = observer.NewObserver[Event]()
)

func main() {
    rspCh, cancelFunc := obsrv.Subscribe()
    defer cancelFunc()
    go func() {
        for {
            fmt.Printf("Received event: %v\n", <-rspCh)
        }
    }()
    fmt.Println("Registered Clients: ", obsrv.Clients())

    obsrv.NotifyAll(Event{
        ID:      i,
        Message: "Hello World",
    })
}
```

If you would like to see a more detailed example, please take a look at the [observer](_example/main.go) example.
