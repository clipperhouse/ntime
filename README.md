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

// Time passes...
time.Sleep(100 * time.Millisecond)

now := ntime.Now()
age := now.Sub(thing.Time)

if age > 100*time.Millisecond {
    //... do a thing with the thing
}

```

[![Tests](https://github.com/clipperhouse/ntime/actions/workflows/tests.yml/badge.svg)](https://github.com/clipperhouse/ntime/actions/workflows/tests.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/clipperhouse/ntime.svg)](https://pkg.go.dev/github.com/clipperhouse/ntime)

### Motivation

The Go stdlib `time` package offers [monotonic time](https://pkg.go.dev/time#hdr-Monotonic_Clocks), which is a nice bit of robustness against a changing system clock. You should usually just use that.

If you're an optimizer like me, and are [building something](https://github.com/clipperhouse/rate) that primarily cares about _relative_ time between two operations, then `ntime.Time` offers an 8-byte type, vs the 24-byte `time.Time`, while still offering monotonicity.

ntime also offers many of the convenience methods from the stdlib `time`, such as `Sub`, `After`, etc.

`ntime.Time` is a relative time against an arbitrary constant epoch, and is meant only for comparisons. Do not mistake it for system ("wall") time.

### Implementation

I naïvely began with the simple idea of using `Nanoseconds()` / `UnixNano()` from the stdlib `time` package. Those are integers!

But then I had the good sense to wonder if they are monotonic in the way that `time.Now()` is. Turns out [they are not](https://chatgpt.com/share/689f6a5d-2f64-8007-b1cc-3bdf10cfee20).

A changing system clock is unlikely, until it isn't, and I would like to eliminate that class of surprise.

One can get both (monotonic + integer) by capturing an epoch via `time.Now()` at system start, and then using `time.Since()`. This package makes that convenient.

### Performance

`ntime` is this package and `time` is the standard library. Looks like `ntime.Now()` is about 2x faster than `time.Now()` on my machine. Other operations are about the same. See BENCHMARKS.txt.

```
goos: darwin
goarch: arm64
pkg: github.com/clipperhouse/ntime
cpu: Apple M2
Benchmark_Now/ntime/Serial-8  	             15.91 ns/op	       0 B/op	       0 allocs/op
Benchmark_Now/ntime/Parallel-8         	     3.334 ns/op	       0 B/op	       0 allocs/op
Benchmark_Now/time/Serial-8            	     32.25 ns/op	       0 B/op	       0 allocs/op
Benchmark_Now/time/Parallel-8          	     10.86 ns/op	       0 B/op	       0 allocs/op
Benchmark_Since/ntime/Serial-8         	     15.90 ns/op	       0 B/op	       0 allocs/op
Benchmark_Since/ntime/Parallel-8       	     3.425 ns/op	       0 B/op	       0 allocs/op
Benchmark_Since/time/Serial-8          	     15.50 ns/op	       0 B/op	       0 allocs/op
Benchmark_Since/time/Parallel-8        	     2.954 ns/op	       0 B/op	       0 allocs/op
Benchmark_After/ntime/Serial-8         	     1.898 ns/op	       0 B/op	       0 allocs/op
Benchmark_After/ntime/Parallel-8       	    0.2988 ns/op	       0 B/op	       0 allocs/op
Benchmark_After/time/Serial-8          	     1.920 ns/op	       0 B/op	       0 allocs/op
Benchmark_After/time/Parallel-8        	    0.3602 ns/op	       0 B/op	       0 allocs/op
Benchmark_Add/ntime/Serial-8           	     1.894 ns/op	       0 B/op	       0 allocs/op
Benchmark_Add/ntime/Parallel-8         	    0.3067 ns/op	       0 B/op	       0 allocs/op
Benchmark_Add/time/Serial-8            	     3.644 ns/op	       0 B/op	       0 allocs/op
Benchmark_Add/time/Parallel-8          	    0.6965 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/ntime/Serial-8           	     1.916 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/ntime/Parallel-8         	    0.3225 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/time/Serial-8            	     1.967 ns/op	       0 B/op	       0 allocs/op
Benchmark_Sub/time/Parallel-8          	    0.3756 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/clipperhouse/ntime
```
