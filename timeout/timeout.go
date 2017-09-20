package timeout

import (
	"time"
)

// SetTimeout exec fn after ms
func SetTimeout(fn func(), d time.Duration) {
	go func() {
		t := time.NewTimer(d)
		<-t.C
		fn()
	}()
}
