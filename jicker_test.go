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
	for gotTime := range c {
		if i >= 3 {
			cancelFunc()
		}
		log.Printf("[debug] %v", gotTime)
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
	for gotTime := range c {
		if i >= 1 {
			cancelFunc()
		}
		log.Printf("[debug] %v", gotTime)
		i++
	}

	end := time.Now()

	elapsedDuration := end.Sub(begin)
	if float64(elapsedDuration) < 5*math.Pow10(8) || float64(elapsedDuration) > 5.2*math.Pow10(8) {
		t.Errorf("elapsed time is out of range unexpectedly; elapsedDuration = %v", float64(elapsedDuration))
	}
}

func TestJicker_TickBetween(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	begin := time.Now()

	i := 1
	c, err := NewJicker().TickBetween(ctx, 400*time.Millisecond, 600*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	for gotTime := range c {
		if i >= 3 {
			cancelFunc()
		}
		log.Printf("[debug] %v", gotTime)
		i++
	}

	end := time.Now()

	elapsedDuration := end.Sub(begin)
	if float64(elapsedDuration) < 1.2*math.Pow10(9) || float64(elapsedDuration) > 1.8*math.Pow10(9) {
		t.Errorf("elapsed time is out of range unexpectedly; elapsedDuration = %v", float64(elapsedDuration))
	}
}

func TestJicker_TickBetween_WithFixed(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	begin := time.Now()

	i := 1
	c, err := NewJicker().TickBetween(ctx, 500*time.Millisecond, 500*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	for gotTime := range c {
		if i >= 1 {
			cancelFunc()
		}
		log.Printf("[debug] %v", gotTime)
		i++
	}

	end := time.Now()

	elapsedDuration := end.Sub(begin)
	if float64(elapsedDuration) < 5*math.Pow10(8) || float64(elapsedDuration) > 5.2*math.Pow10(8) {
		t.Errorf("elapsed time is out of range unexpectedly; elapsedDuration = %v", float64(elapsedDuration))
	}
}

func TestJicker_TickBetween_ShouldRaiseErrorWhenArgumentsAreInverted(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	_, err := NewJicker().TickBetween(ctx, 2*time.Millisecond, 1*time.Millisecond)
	if err == nil {
		t.Fatal("expected error has been nil")
	}
}

func TestJicker_Tick_EnsureNonBlockingEvenIfClientDoesNotConsumeChannel(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	c := NewJicker().Tick(ctx, 10*time.Millisecond, 0)
	time.Sleep(100 * time.Millisecond)

	gotTime := <-c
	log.Printf("[debug] %v", gotTime)
}

func ExampleJicker_Tick() {
	// if this `ctx` has done, ticking stops and it closes the ticker channel.
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// it ticks by jittered interval duration (i.e. 1±5% sec); it evaluates the duration with the jitter factor every time.
	c := NewJicker().Tick(ctx, 1*time.Second, 0.05)
	for t := range c {
		log.Printf("tick: %v", t)
	}
}

func ExampleJicker_TickBetween() {
	// if this `ctx` has done, ticking stops and it closes the ticker channel.
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// it ticks by jittered interval duration; random duration is between [minimumDuration, maximumDuration].
	// it evaluates the duration with the jitter factor every time.
	//
	// in this case, the interval duration would be the random value between 1 sec and 2 secs.
	c, err := NewJicker().TickBetween(ctx, 1*time.Second, 2*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	for t := range c {
		log.Printf("tick: %v", t)
	}
}
