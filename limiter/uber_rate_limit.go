package limiter

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/benbjohnson/clock"
)

type Limiter interface {
	Take() time.Time
}

type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

type config struct {
	clock Clock
	slack int
	per   time.Duration
}

func buildConfig(opts []Option) config {
	c := config{
		clock: clock.New(),
		slack: 10,
		per:   time.Second,
	}
	for _, opt := range opts {
		opt.apply(&c)
	}
	return c
}

type Option interface {
	apply(*config)
}

type state struct {
	last     time.Time
	sleepFor time.Duration
}

type atomicLimiter struct {
	state      unsafe.Pointer
	padding    [56]byte // cache line size - state pointer size = 64 -8 ; created to avoid false sharing
	perRequest time.Duration
	maxSlack   time.Duration
	clock      Clock
}

func newAtomicBased(rate int, opts ...Option) *atomicLimiter {
	config := buildConfig(opts)
	perRequest := config.per / time.Duration(rate)
	l := &atomicLimiter{
		perRequest: perRequest,
		maxSlack:   -1 * time.Duration(config.slack) * perRequest,
		clock:      config.clock,
	}
	initialState := state{
		last:     time.Time{},
		sleepFor: 0,
	}
	atomic.StorePointer(&l.state, unsafe.Pointer(&initialState))
	return l
}

func (t *atomicLimiter) Take() time.Time {
	var (
		newState state
		token    bool
		interval time.Duration
	)
	for !token {
		now := t.clock.Now()

		previousStatePointer := atomic.LoadPointer(&t.state)
		oldState := (*state)(previousStatePointer)
		newState = state{
			last:     now,
			sleepFor: oldState.sleepFor,
		}
		if oldState.last.IsZero() {
			token = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
			continue
		}
		newState.sleepFor += t.perRequest - now.Sub(oldState.last)
		if newState.sleepFor < t.maxSlack {
			newState.sleepFor = t.maxSlack
		}
		if newState.sleepFor > 0 {
			newState.last = newState.last.Add(newState.sleepFor)
			interval, newState.sleepFor = newState.sleepFor, 0
		}
		token = atomic.CompareAndSwapPointer(&t.state, previousStatePointer, unsafe.Pointer(&newState))
	}
	t.clock.Sleep(interval)
	return newState.last
}
