package channel

import (
	"sync"
	"time"
)

// Channel define channel interface
type Channel interface {
	Done()
	IsDone() bool
	Send(...interface{}) bool
	Get(deadline time.Duration) ([]interface{}, bool)
}

type wrapCh struct {
	sync.RWMutex

	_ch    chan interface{}
	isDone bool
}

// New channel
func New(retCnt int) Channel {
	return &wrapCh{_ch: make(chan interface{}, retCnt)}
}

func (ch *wrapCh) Done() {
	ch.Lock()
	defer ch.Unlock()

	if ch.isDone {
		return
	}

	close(ch._ch)
	ch.isDone = true
}

func (ch *wrapCh) IsDone() bool {
	ch.RLock()
	defer ch.RUnlock()

	if ch.isDone == true {
		return true
	}

	return false
}

func (ch *wrapCh) Send(vals ...interface{}) bool {
	ch.RLock()
	defer ch.RUnlock()

	if ch.isDone == true {
		return false
	}

	for _, v := range vals {
		if len(ch._ch) >= cap(ch._ch) {
			return false
		}

		ch._ch <- v
	}

	return true
}

func (ch *wrapCh) Get(deadline time.Duration) ([]interface{}, bool) {
	ch.RLock()
	defer ch.RUnlock()

	var ret []interface{}
	cnt, chCap := 0, cap(ch._ch)

	end := time.Now().Add(deadline)

	for {
		if end.Sub(time.Now()) <= 0 {
			return nil, false
		}

		if cnt >= chCap {
			return ret, true
		}

		select {
		case v, ok := <-ch._ch:
			if ok {
				ret = append(ret, v)
				cnt++
			}
		default:
		}
	}
}
