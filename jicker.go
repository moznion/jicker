package jicker

import (
	"context"
	"math/rand"
	"time"
)

// Jicker is a struct for jittered-ticker.
type Jicker struct {
	rng *rand.Rand
}

// NewJicker makes a new Jicker instance with new random number generator.
func NewJicker() *Jicker {
	return &Jicker{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Tick starts a ticker with jittering.
// This function returns a "Time" channel that is notified on ticking.
//
// parameters:
// - ctx: A context object. If this context has done, ticker stops and it closes the response channel.
// - duration: Base duration for ticking.
// - jitterFactor: A factor to jitter the ticking duration. For example, if this value is 0.05, each ticker duration jittered ±5%.
func (ji *Jicker) Tick(ctx context.Context, duration time.Duration, jitterFactor float64) <-chan time.Time {
	if jitterFactor <= 0 {
		jitterFactor = 0
	}

	timeCh := make(chan time.Time)

	go func() {
		defer close(timeCh)

		t := time.Now()
		sleepingCh := make(chan struct{})

		for {
			jitteredDuration := time.Duration(ji.jitter(float64(duration), jitterFactor))
			go func() {
				time.Sleep(jitteredDuration)
				t = t.Add(jitteredDuration)
				sleepingCh <- struct{}{}
			}()

			select {
			case <-sleepingCh:
				select {
				case timeCh <- t:
				default:
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return timeCh
}

func (ji *Jicker) jitter(duration float64, jitterFactor float64) float64 {
	jitterDelta := duration * jitterFactor
	jitterMin := duration - jitterDelta
	jitterMax := duration + jitterDelta

	// Get a random value from the range [minInterval, maxInterval].
	// The formula used below has a +1 because if the minInterval is 1 and the maxInterval is 3 then
	// we want a 33% chance for selecting either 1, 2 or 3.
	//
	// see also: https://github.com/cenkalti/backoff/blob/c2975ffa541a1caeca5f76c396cb8c3e7b3bb5f8/exponential.go#L154-L157
	return jitterMin + ji.rng.Float64()*(jitterMax-jitterMin+1)
}
