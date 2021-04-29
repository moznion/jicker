package jicker

import (
	"context"
	"log"
	"math"
	"testing"
	"time"
)

func TestJicker_Tick_WithJitter(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	begin := time.Now()

	i := 1
	c := NewJicker().Tick(ctx, 500*time.Millisecond, 0.2)
	for t := range c {
		if i >= 3 {
			cancelFunc()
		}
		log.Printf("[debug] %v", t)
		i++
	}

	end := time.Now()

	elapsedDuration := end.Sub(begin)
	if float64(elapsedDuration) < 1.2*math.Pow10(9) || float64(elapsedDuration) > 1.8*math.Pow10(9) {
		t.Errorf("elapsed time is out of range unexpectedly; elapsedDuration = %v", float64(elapsedDuration))
	}
}

func TestJicker_Tick_WithFixed(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	begin := time.Now()

	i := 1
	c := NewJicker().Tick(ctx, 500*time.Millisecond, 0)
	for t := range c {
		if i >= 3 {
			cancelFunc()
		}
		log.Printf("[debug] %v", t)
		i++
	}

	end := time.Now()

	elapsedDuration := end.Sub(begin)
	if float64(elapsedDuration) < 1.5*math.Pow10(9) || float64(elapsedDuration) > 1.52*math.Pow10(9) {
		t.Errorf("elapsed time is out of range unexpectedly; elapsedDuration = %v", float64(elapsedDuration))
	}
}

func ExampleJicker_Tick() {
	// if this `ctx` has done, ticking stops and it closes the ticker channel.
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// it ticks by jittered duration (i.e. 1±5% sec); it evaluates the duration with the jitter factor every time.
	c := NewJicker().Tick(ctx, 1*time.Second, 0.05)
	for t := range c {
		log.Printf("tick: %v", t)
	}
}
