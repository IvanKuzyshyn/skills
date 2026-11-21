package ratelimit

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time          { return c.t }
func (c *fakeClock) advance(d time.Duration) { c.t = c.t.Add(d) }

func TestAllow(t *testing.T) {
	tests := []struct {
		name    string
		rate    float64
		burst   int
		elapsed time.Duration
		take    []int
		want    []bool
	}{
		{"full bucket", 1, 3, 0, []int{1, 1, 1, 1}, []bool{true, true, true, false}},
		{"refills over time", 2, 2, 500 * time.Millisecond, []int{2, 1, 1}, []bool{true, true, false}},
		{"never exceeds burst", 100, 2, time.Hour, []int{2, 3}, []bool{true, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clk := &fakeClock{t: time.Unix(0, 0)}
			b := newWithClock(tt.rate, tt.burst, clk.now)
			for i, n := range tt.take {
				if i == 1 {
					clk.advance(tt.elapsed)
				}
				if got := b.Allow(n); got != tt.want[i] {
					t.Errorf("Allow(%d) call %d = %v, want %v", n, i, got, tt.want[i])
				}
			}
		})
	}
}

func TestWaitExceedsBurst(t *testing.T) {
	b := New(1, 1)
	if err := b.Wait(context.Background(), 2); !errors.Is(err, ErrExceedsBurst) {
		t.Fatalf("Wait(2) error = %v, want %v", err, ErrExceedsBurst)
	}
}

func TestWaitCancelled(t *testing.T) {
	b := New(0.001, 1)
	b.Allow(1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := b.Wait(ctx, 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("Wait error = %v, want %v", err, context.Canceled)
	}
}

func TestConcurrentAllow(t *testing.T) {
	b := New(0, 50)
	var wg sync.WaitGroup
	var mu sync.Mutex
	allowed := 0
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if b.Allow(1) {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if allowed != 50 {
		t.Errorf("allowed = %d, want 50", allowed)
	}
}
