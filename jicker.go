package jicker

import (
	"context"
	"errors"
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

// Tick starts a ticker with jittering based on factor.
// This function returns a "Time" channel that is notified on ticking.
//
// parameters:
// - ctx: A context object. If this context has done, ticker stops and it closes the response channel.
// - baseIntervalDuration: Base interval duration for ticking.
// - jitterFactor: A factor to jitter the interval duration. For example, if this value is 0.05, each interval duration jittered ±5%.
func (ji *Jicker) Tick(ctx context.Context, baseIntervalDuration time.Duration, jitterFactor float64) <-chan time.Time {
	if jitterFactor <= 0 {
		jitterFactor = 0
	}

	return ji.tick(ctx, newFactoredJitter(baseIntervalDuration, jitterFactor, ji.rng))
}

// TickBetween starts a ticker with jittered random value between [minimumIntervalDuration, maximumIntervalDuration].
// This function returns a "Time" channel that is notified on ticking.
//
// parameters:
// - ctx: A context object. If this context has done, ticker stops and it closes the response channel.
// - minimumIntervalDuration: The minimum interval duration for jittering. This must be smaller than maximumIntervalDuration.
// - maximumIntervalDuration: The maximum interval duration for jittering. This must be larger than minimumIntervalDuration.
func (ji *Jicker) TickBetween(ctx context.Context, minimumIntervalDuration time.Duration, maximumIntervalDuration time.Duration) (<-chan time.Time, error) {
	if maximumIntervalDuration < minimumIntervalDuration {
		return nil, errors.New("minimumIntervalDuration must be smaller than maximumIntervalDuration, but it's not")
	}

	return ji.tick(ctx, newRangeJitter(minimumIntervalDuration, maximumIntervalDuration, ji.rng)), nil
}

func (ji *Jicker) tick(ctx context.Context, jitterable jitterable) <-chan time.Time {
	timeCh := make(chan time.Time)

	go func() {
		defer close(timeCh)

		for {
			jitteredDuration := time.Duration(jitterable.jitter())

			select {
			case t := <-time.After(jitteredDuration):
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

type jitterable interface {
	jitter() float64
}

type factoredJitter struct {
	baseIntervalDuration float64
	jitterFactor         float64
	rng                  *rand.Rand
}

func newFactoredJitter(baseIntervalDuration time.Duration, jitterFactor float64, rng *rand.Rand) *factoredJitter {
	return &factoredJitter{
		baseIntervalDuration: float64(baseIntervalDuration),
		jitterFactor:         jitterFactor,
		rng:                  rng,
	}
}

func (fj *factoredJitter) jitter() float64 {
	jitterDelta := fj.baseIntervalDuration * fj.jitterFactor
	jitterMin := fj.baseIntervalDuration - jitterDelta
	jitterMax := fj.baseIntervalDuration + jitterDelta
	return jitter(jitterMin, jitterMax, fj.rng)
}

type rangeJitter struct {
	minimumIntervalDuration float64
	maximumIntervalDuration float64
	rng                     *rand.Rand
}

func newRangeJitter(minimumIntervalDuration time.Duration, maximumIntervalDuration time.Duration, rng *rand.Rand) *rangeJitter {
	return &rangeJitter{
		minimumIntervalDuration: float64(minimumIntervalDuration),
		maximumIntervalDuration: float64(maximumIntervalDuration),
		rng:                     rng,
	}
}

func (rj *rangeJitter) jitter() float64 {
	return jitter(rj.minimumIntervalDuration, rj.maximumIntervalDuration, rj.rng)
}

func jitter(min float64, max float64, rng *rand.Rand) float64 {
	// Get a random value from the range [minInterval, maxInterval].
	// The formula used below has a +1 because if the minInterval is 1 and the maxInterval is 3 then
	// we want a 33% chance for selecting either 1, 2 or 3.
	//
	// see also: https://github.com/cenkalti/backoff/blob/c2975ffa541a1caeca5f76c396cb8c3e7b3bb5f8/exponential.go#L154-L157
	return min + rng.Float64()*(max-min+1)
}
