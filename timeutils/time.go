package timeutils

import (
	"math"
	"time"
)

const (
	Format       = "2006-01-02 15:04:05"
	RFC3339      = "2006-01-02T15:04:05Z07:00"
	RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"
	RFC3339Micro = "2006-01-02T15:04:05.999999Z07:00"
	RFC3339Nano  = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen      = "3:04PM"
	KitchenSec   = "3:04:05PM"
	Date         = "2006-01-02"
	Date2        = "2006/01/02"
	Stamp        = "15:04:05"
	StampMilli   = "15:04:05.000"
	StampMicro   = "15:04:05.000000"
	StampNano    = "15:04:05.000000000"
)

var (
	// CSTZone 中国标准时间
	CSTZone = time.FixedZone("CST", 8*3600)
)

// Min returns the min duration between x and y.
func Min(x, y time.Duration) time.Duration {
	if x < math.MinInt64 || y < math.MinInt64 {
		return math.MinInt64
	}
	if x < y {
		return x
	}
	return y
}

// Max returns the max duration between x and y.
func Max(x, y time.Duration) time.Duration {
	if x > math.MaxInt64 || y > math.MaxInt64 {
		return math.MaxInt64
	}
	if x > y {
		return x
	}
	return y
}

func Add(x, y time.Duration) time.Duration {
	return x + y
}
