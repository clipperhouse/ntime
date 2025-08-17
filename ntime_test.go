package ntime

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// @clipperhouse: I don't know if we can actually test monotonicity
// against system clock changes. These tests are something I suppose.
// https://chatgpt.com/share/689f6a5d-2f64-8007-b1cc-3bdf10cfee20

func TestNTime_Now_Vs_TimeNow_Serial(t *testing.T) {
	t.Parallel()

	// Test that consecutive calls to ntime.Now() return non-decreasing
	// values. I believe that returing the same time is acceptable
	// for a monotonic clock, for our purposes.
	{
		prev := Now()
		for i := range 10000 {
			current := Now()
			require.GreaterOrEqual(t, current, prev,
				"Time %d (%d) should not be less than time %d (%d)",
				i, current, i-1, prev)
			prev = current
		}
	}

	// For comparison, test that consecutive calls to time.Now()
	// do the same.
	{
		prev := time.Now()
		for i := range 10000 {
			current := time.Now()
			require.GreaterOrEqual(t, current, prev,
				"Time %d (%v) should not be less than time %d (%v)",
				i, current, i-1, prev)
			prev = current
		}
	}
}

func TestNtime_Now_Vs_TimeNow_Concurrent(t *testing.T) {
	t.Skip("skipping because it's not a good test of monotonicity, and is non-deterministic")

	// Test the extent to which ntime.Now() offers monotonicity
	// vs system time.Now()

	// @clipperhouse: I am not convinced this is a good test of
	// monotonicity. Both ntime.Now() and time.Now() show
	// "violations", but I suspect that's an artifact of
	// goroutines and channels, and not the time calculations.
	// Consider this test an exploration.

	t.Parallel()

	const concurrency = 10
	const calls = 100
	const total = concurrency * calls

	// Test ntime.Now() concurrent monotonicity (this package)
	nviolations := 0
	{
		c := make(chan Time, concurrency*calls)
		var wg sync.WaitGroup
		for range concurrency {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range calls {
					c <- Now()
				}
			}()
		}
		wg.Wait()
		close(c)

		prev := <-c
		for ntime := range c {
			if ntime < prev {
				nviolations++
			}
			prev = ntime
		}

		pct := 100 * nviolations / (total - 1)
		t.Logf("ntime.Now() violations: %d/%d (%d%%)",
			nviolations, total-1, pct,
		)
	}

	// Test time.Now() concurrent monotonicity (stdlib)
	tviolations := 0
	{
		c := make(chan time.Time, concurrency*calls)
		var wg sync.WaitGroup
		for range concurrency {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range calls {
					c <- time.Now()
				}
			}()
		}
		wg.Wait()
		close(c)

		prev := <-c
		for ttime := range c {
			if ttime.Before(prev) {
				tviolations++
			}
			prev = ttime
		}

		pct := 100 * tviolations / (total - 1)
		t.Logf("time.Now() violations: %d/%d (%d%%)",
			tviolations, total-1, pct,
		)
	}

	// Report if they greatly differ, *maybe* it's a real
	// regression.
	require.Less(t, nviolations, 2*tviolations,
		"ntime.Now() has 2x more monotonicity violations (%d) than time.Now() (%d)",
		nviolations, tviolations)
	require.Less(t, tviolations, 2*nviolations,
		"time.Now() has 2x more monotonicity violations (%d) than ntime.Now() (%d)",
		tviolations, nviolations)
}

func TestTime_Add(t *testing.T) {
	t.Parallel()

	now := Now()
	duration := time.Second

	result := now.Add(duration)
	expected := now + Time(duration.Nanoseconds())

	require.Equal(t, expected, result, "Add: got %d, want %d", result, expected)
}

func TestTime_Sub(t *testing.T) {
	t.Parallel()

	now := Now()
	time.Sleep(time.Millisecond)
	later := Now()

	duration := later.Sub(now)

	require.Greater(t, duration, time.Duration(0), "Sub: duration should be positive, got %v", duration)

	// Verify the duration is reasonable (should be around 1ms)
	require.GreaterOrEqual(t, duration, time.Microsecond, "Sub: duration too small: %v", duration)
	require.LessOrEqual(t, duration, time.Second, "Sub: duration too large: %v", duration)
}

func TestTime_After(t *testing.T) {
	t.Parallel()

	now := Now()
	time.Sleep(time.Millisecond)
	later := Now()

	require.True(t, later.After(now), "After: later time should be after earlier time")
	require.False(t, now.After(later), "After: earlier time should not be after later time")
}

func TestTime_Before(t *testing.T) {
	t.Parallel()

	now := Now()
	time.Sleep(time.Millisecond)
	later := Now()

	require.True(t, now.Before(later), "Before: earlier time should be before later time")
	require.False(t, later.Before(now), "Before: later time should not be before earlier time")
}

func BenchmarkNow(b *testing.B) {
	for b.Loop() {
		Now()
	}
}

func BenchmarkNow_Parallel(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Now()
		}
	})
}

func BenchmarkTime_Add(b *testing.B) {
	now := Now()
	duration := time.Second

	for b.Loop() {
		now.Add(duration)
	}
}

func BenchmarkTime_Sub(b *testing.B) {
	now := Now()
	time.Sleep(time.Microsecond)
	later := Now()

	for b.Loop() {
		later.Sub(now)
	}
}
