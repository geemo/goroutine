package race

import (
	"errors"
	"test/goroutine/race/channel"
	"time"
)

// Race promise.race
func Race(retCnt int, deadline time.Duration, fns ...func(ch channel.Channel)) ([]interface{}, error) {
	ch := channel.New(retCnt)

	for _, fn := range fns {
		go fn(ch)
	}

	vals, ok := ch.Get(deadline)

	if !ok {
		return nil, errors.New("exec failed")
	}

	return vals, nil
}
