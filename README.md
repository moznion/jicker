# jicker [![Go](https://github.com/moznion/jicker/actions/workflows/check.yml/badge.svg)](https://github.com/moznion/jicker/actions/workflows/check.yml)

A jittered-ticker library for go.

[![Go Reference](https://pkg.go.dev/badge/github.com/moznion/jicker.svg)](https://pkg.go.dev/github.com/moznion/jicker)

## Usage

```go
package jicker

import (
	"context"
	"log"
	"math"
	"testing"
	"time"
)

func main() {
    // if this `ctx` has done, ticking stops and it closes the ticker channel.
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// it ticks by jittered duration (i.e. 1±5% sec)
	c := NewJicker().Tick(ctx, 1*time.Second, 0.05)
	for t := range c {
		log.Printf("tick: %v", t)
	}
}
```

## See also

- https://golang.org/pkg/time/#Tick

## Author

moznion (<moznion@mail.moznion.net>)

