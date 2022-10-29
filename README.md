# observer

[![tests](https://github.com/leonsteinhaeuser/observer/actions/workflows/tests.yml/badge.svg)](https://github.com/leonsteinhaeuser/observer/actions/workflows/tests.yml)
[![codecov](https://codecov.io/gh/leonsteinhaeuser/observer/branch/main/graph/badge.svg?token=WQUVF65EVF)](https://codecov.io/gh/leonsteinhaeuser/observer)
[![release date](https://img.shields.io/github/release-date/leonsteinhaeuser/observer)](https://img.shields.io/github/release-date/leonsteinhaeuser/observer)
[![commits since release](https://img.shields.io/github/commits-since/leonsteinhaeuser/observer/latest)](https://img.shields.io/github/commits-since/leonsteinhaeuser/observer/latest)
[![open: bugs](https://img.shields.io/github/issues/leonsteinhaeuser/observer/bug)](https://img.shields.io/github/issues/leonsteinhaeuser/observer/bug)
[![open: feature requests](https://img.shields.io/github/issues/leonsteinhaeuser/observer/feature%20request)](https://img.shields.io/github/issues/leonsteinhaeuser/observer/feature%20request)
[![issues closed](https://img.shields.io/github/issues-closed/leonsteinhaeuser/observer)](https://img.shields.io/github/issues-closed/leonsteinhaeuser/observer)
[![license](https://img.shields.io/github/license/leonsteinhaeuser/observer)](https://img.shields.io/github/license/leonsteinhaeuser/observer)

This library implements the pub/sub pattern in a generic way. It uses Go's generic types to declare the type of the event.

## Differences between v1 and v2

The v2 release contains breaking changes. The most important ones are:

- The `Observable` interface does not provide a `Clients() int64` method anymore.
- The constructor `NewObserver()` has been removed. Instead, use `new(observer.Observer[T])`.
- The `Observer` has become a `NotifyTimeout` that can be used to set a timeout for the `NotifyAll` method. The default value is `5 * time.Second`.

## Usage

```bash
go get github.com/leonsteinhaeuser/observer/v2
```

## Example

```go
package main

import (
    "fmt"
    "github.com/leonsteinhaeuser/observer/v2"
)

type Event struct {
    ID      int
    Message string
}

var (
    obsrv *observer.Observer[Event] = new(observer.Observer[Event])
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

If you would like to see a more detailed example, please take a look at the [observer](_example/handler/main.go) example.
