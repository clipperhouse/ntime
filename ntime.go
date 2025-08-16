// Package ntime provides a monotonic time, expressed as an
// int64 nanosecond count since an arbitrary epoch. Intended
// for applications which only require relative time measurements.
//
// The epoch is time.Now() at package initialization.
package ntime

import (
	"time"
)

// Time is an int64 which represents a monotonic time source,
// expressed as a nanosecond count since an arbitrary static
// epoch.
//
// Intended as an optimization to store relative time as an
// int64 (8 bytes), instead of a time.Time (24 bytes).
//
// See https://chatgpt.com/share/689f6a5d-2f64-8007-b1cc-3bdf10cfee20
type Time int64

var epoch = time.Now()

// Epoch returns the static epoch time that is used as the
// basis for ntime.Time calculations. The epoch is time.Now()
// at package initialization.
//
// By itself, it's probably not useful, but is offered for
// diagnostic purposes.
func Epoch() time.Time {
	return epoch
}

// Now returns the current relative monotonic time since
// an arbitrary static epoch. Intended for use similar to
// time.Now(), but for applications which only require
// relative time measurements.
func Now() Time {
	return Time(time.Since(epoch).Nanoseconds())
}

// ToTime adds monotonic int64 nanoseconds (ntime.Time) to
// the static epoch, returning a time.Time that will usually
// be close to time.Now().
//
// ⚠️ This might or might not match time.Now(), and
// in the presence of system clock changes, it might be
// surprising. Use sparingly.
//
// Intended as a shim for use with clipperhouse/rate.
func (t Time) ToTime() time.Time {
	return epoch.Add(time.Duration(t))
}

// various shims to look like time.Time methods

func (t Time) Add(d time.Duration) Time {
	return t + Time(d.Nanoseconds())
}

func (t Time) Sub(u Time) time.Duration {
	return time.Duration(t - u)
}

func (t Time) After(u Time) bool {
	return t > u
}

func (t Time) Before(u Time) bool {
	return t < u
}

func (t Time) BeforeOrEqual(u Time) bool {
	return t <= u
}

func (t Time) AfterOrEqual(u Time) bool {
	return t >= u
}
