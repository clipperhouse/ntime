### ntime

ntime provides a monotonic time, expressed as an `int64` nanosecond count. It is intended for applications which only require relative time measurements, and wish to optimize memory usage and (likely) speed.

```go
go get github.com/clipperhouse/ntime
```

```go
import "github.com/clipperhouse/ntime"

type MyThing struct {
	Time ntime.Time
}

thing := MyThing{
    Time: ntime.Now(),
}

// time passes...

now := ntime.Now()
age := now.Sub(thing.Time)

if age > time.Duration(30*time.Second) {
    //... do a thing to the thing
}
```

[![Tests](https://github.com/clipperhouse/ntime/actions/workflows/tests.yml/badge.svg)](https://github.com/clipperhouse/ntime/actions/workflows/tests.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/clipperhouse/ntime.svg)](https://pkg.go.dev/github.com/clipperhouse/ntime)

### Motivation

The Go stdlib `time` package offers [monotonic time](https://pkg.go.dev/time#hdr-Monotonic_Clocks), which is a nice bit of robustness against a changing system clock. You should usually just use that.

If you're an optimizer like me, and are [building something](https://github.com/clipperhouse/rate) that primarily cares about _relative_ time between two operations, then `ntime.Time` offers an 8-byte type, vs the 24-byte `time.Time`.

ntime also offers many of the convenience methods from the stdlib `time`, such as `Sub`, `After`, etc.

`ntime.Time` is a relative time against an arbitrary static epoch, and is meant only for comparisons. Do not mistake it for system ("wall") time.

### Implementation

I naïvely began with the simple idea of using `Nanoseconds()` / `UnixNano()` from the stdlib `time` package. Those are integers!

But then I had the good sense to wonder if they are monotonic in the way that `time.Now()` is. Turns out [they are not](https://chatgpt.com/share/689f6a5d-2f64-8007-b1cc-3bdf10cfee20).

A changing systems clock is a class of bug that is unlikely, until it isn't, and I would like to eliminate that class of surprise.

One can get both (monotonic + integer) by capturing an epoch via `time.Now()` at system start, and then using `time.Since()`. This package makes that convenient.
