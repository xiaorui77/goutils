package wait

import "time"

// WaitUntil will wait (block) until the function status is true by checks every 1000ms.
func WaitUntil(f func() bool) {
	waitUntil(f, time.Millisecond*1000)
}

// waitUntil will wait (block) until the function status is true by checks every t interval.
func waitUntil(f func() bool, t time.Duration) {
	// The minimum interval is 100 ms
	if t < time.Millisecond*100 {
		t = time.Millisecond * 100
	}

	for {
		if f() {
			return
		}
		time.Sleep(t)
	}
}
