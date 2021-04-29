# jicker [![Go](https://github.com/moznion/jicker/actions/workflows/check.yml/badge.svg)](https://github.com/moznion/jicker/actions/workflows/check.yml)

A jittered-ticker library for go.

[![Go Reference](https://pkg.go.dev/badge/github.com/moznion/jicker.svg)](https://pkg.go.dev/github.com/moznion/jicker)

## Usage

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/moznion/jicker"
)

func main() {
	// if this `ctx` has done, ticking stops and it closes the ticker channel.
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// it ticks by jittered duration (i.e. 1±5% sec); it evaluates the duration with the jitter factor every time.
	c := jicker.NewJicker().Tick(ctx, 1*time.Second, 0.05)
	for t := range c {
		log.Printf("tick: %v", t)
	}
}
```

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/moznion/jicker"
)

func main() {
	// if this `ctx` has done, ticking stops and it closes the ticker channel.
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// it ticks by jittered interval duration; random duration is between [minimumDuration, maximumDuration].
	// it evaluates the duration with the jitter factor every time.
	//
	// in this case, the interval duration would be the random value between 1 sec and 2 secs.
	c, err := jicker.NewJicker().TickBetween(ctx, 1*time.Second, 2*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	for t := range c {
		log.Printf("tick: %v", t)
	}
}
```

## See also

- https://golang.org/pkg/time/#Tick

## Author

moznion (<moznion@mail.moznion.net>)

