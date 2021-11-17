package math

import (
	"math"
	"time"
)

// MinDuration returns the min duration between x and y.
func MinDuration(x, y time.Duration) time.Duration {
	if x < math.MinInt64 || y < math.MinInt64 {
		return math.MinInt64
	}
	if x < y {
		return x
	}
	return y
}

// MaxDuration returns the max duration between x and y.
func MaxDuration(x, y time.Duration) time.Duration {
	if x > math.MaxInt64 || y > math.MaxInt64 {
		return math.MaxInt64
	}
	if x > y {
		return x
	}
	return y
}
