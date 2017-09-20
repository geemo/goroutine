package eachLimit

import (
	"errors"
)

// EachLimit each do fn
func EachLimit(
	arr []interface{},
	limit int,
	fn func(interface{}, func(error)),
	done func(error)) {

	arrLen := len(arr)
	if arrLen == 0 {
		done(nil)
		return
	}

	if limit <= 0 || limit > arrLen {
		done(errors.New("error limit"))
		return
	}

	_ch := make(chan error, arrLen)
	_next := func(err error) {
		if err != nil {
			done(err)
			return
		}

		_ch <- nil
	}

	go func() {
		idx := 0
		for {
			if idx >= arrLen {
				done(nil)
				return
			}

			realLimit := 0
			for i := 0; i < limit && idx < arrLen; i++ {
				realLimit++
				go fn(arr[idx], _next)
				idx++
			}

			for realLimit > 0 {
				select {
				case err, ok := <-_ch:
					if ok {
						if err != nil {
							done(err)
							return
						}

						realLimit--
					}
				}
			}
		}

		done(nil)
	}()

}
