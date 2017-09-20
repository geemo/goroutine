package parallel

import (
	"sync"
)

type wrapMap struct {
	sync.RWMutex

	m   map[int]interface{}
	cnt int
}

// Parallel run fns in parallel
func Parallel(fns []func(cb func(error, ...interface{})), done func([]interface{}, error)) {
	_map := &wrapMap{m: make(map[int]interface{})}

	_cb := func(idx int, length int) func(error, ...interface{}) {
		return func(err error, vals ...interface{}) {
			if err != nil {
				done(nil, err)
				return
			}

			_map.Lock()
			defer _map.Unlock()

			_map.m[idx] = vals
			_map.cnt++

			if _map.cnt == length {
				var result []interface{}

				for i := 0; i < length; i++ {
					result = append(result, _map.m[i])
				}

				done(result, nil)
			}
		}
	}

	for idx, fn := range fns {
		go fn(_cb(idx, len(fns)))
	}
}
