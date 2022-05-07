package time

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

const (
	Nanosecond  time.Duration = 1
	Microsecond               = 1000 * Nanosecond
	Millisecond               = 1000 * Microsecond
	Centisecond               = 10 * Millisecond
	Decisecond                = 10 * Centisecond
	Second                    = 10 * Decisecond
	Minute                    = 60 * Second
	Hour                      = 60 * Minute
)

var (
	// CSTZone 中国标准时间
	CSTZone = time.FixedZone("CST", 8*3600)
	// UTCZone 标准时间
	UTCZone = time.UTC
	// ESTZone 美国东部时间
	ESTZone = time.FixedZone("EST", -5*3600)
	// PSTZone 美国太平洋时间
	PSTZone = time.FixedZone("PST", -8*3600)
)

// Min returns the min Duration between x and y.
func Min(x, y time.Duration) time.Duration {
	if x < math.MinInt64 || y < math.MinInt64 {
		return math.MinInt64
	}
	if x < y {
		return x
	}
	return y
}

// Max returns the max Duration between x and y.
func Max(x, y time.Duration) time.Duration {
	if x > math.MaxInt64 || y > math.MaxInt64 {
		return math.MaxInt64
	}
	if x > y {
		return x
	}
	return y
}

// Add returns the Duration x+y.
func Add(x, y time.Duration) time.Duration {
	return x + y
}

// Duration wrap time.Duration
type Duration time.Duration

// DeciSecond return the duration in deci-seconds.
// use: fmt.Printf("%0.1fs", d.DeciSecond())
func (d Duration) DeciSecond() float64 {
	return time.Duration(d).Truncate(Decisecond).Seconds()
}

// CentiSecond return the duration in centi-seconds.
// use: fmt.Printf("%0.2fs", d.CentiSecond())
func (d Duration) CentiSecond() float64 {
	return time.Duration(d).Truncate(Centisecond).Seconds()
}
