package interval

import (
	"time"
)

// SetInterval exec fn pre d time.Duration
func SetInterval(fn func(), d time.Duration) {
	go func() {
		t := time.NewTicker(d)
		for _ = range t.C {
			fn()
		}
	}()
}
