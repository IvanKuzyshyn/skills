package ratelimit

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrExceedsBurst is returned by Wait when n is larger than the bucket can
// ever hold, so waiting would never succeed.
var ErrExceedsBurst = errors.New("ratelimit: request exceeds burst size")

// Bucket is a token bucket. Tokens are added continuously at Rate per second
// up to Burst. A Bucket must be created with New.
type Bucket struct {
	mu     sync.Mutex
	rate   float64
	burst  float64
	tokens float64
	last   time.Time
	now    func() time.Time
}

// New returns a full bucket that refills at rate tokens per second and holds
// at most burst tokens.
func New(rate float64, burst int) *Bucket {
	return newWithClock(rate, burst, time.Now)
}

func newWithClock(rate float64, burst int, now func() time.Time) *Bucket {
	return &Bucket{
		rate:   rate,
		burst:  float64(burst),
		tokens: float64(burst),
		last:   now(),
		now:    now,
	}
}

// Allow reports whether n tokens are available now and, if so, consumes them.
func (b *Bucket) Allow(n int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill()
	if b.tokens < float64(n) {
		return false
	}
	b.tokens -= float64(n)
	return true
}

// Wait blocks until n tokens are available or ctx is done.
func (b *Bucket) Wait(ctx context.Context, n int) error {
	if float64(n) > b.burst {
		return ErrExceedsBurst
	}
	for {
		d, ok := b.reserve(n)
		if ok {
			return nil
		}
		t := time.NewTimer(d)
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
		}
	}
}

// reserve consumes n tokens if available; otherwise it returns how long until
// they will be.
func (b *Bucket) reserve(n int) (time.Duration, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill()
	if b.tokens >= float64(n) {
		b.tokens -= float64(n)
		return 0, true
	}
	missing := float64(n) - b.tokens
	return time.Duration(missing / b.rate * float64(time.Second)), false
}

// refill must be called with b.mu held.
func (b *Bucket) refill() {
	now := b.now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now
	b.tokens = min(b.burst, b.tokens+elapsed*b.rate)
}
