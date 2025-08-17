package ntime_test

import (
	"fmt"
	"time"

	"github.com/clipperhouse/ntime"
)

// Example demonstrates basic usage of ntime for relative time
// measurements.
func Example() {
	type MyThing struct {
		Time ntime.Time
	}

	thing := MyThing{
		Time: ntime.Now(),
	}

	// Time passes...
	time.Sleep(100 * time.Millisecond)

	// Get current time and calculate age
	now := ntime.Now()
	age := now.Sub(thing.Time)

	fmt.Printf("Age > 50ms: %t\n", age > 50*time.Millisecond)
	fmt.Printf("Age > 1 second: %t\n", age > time.Second)

	// Output:
	// Age > 50ms: true
	// Age > 1 second: false
}
