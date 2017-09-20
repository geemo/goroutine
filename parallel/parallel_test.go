package parallel

import (
	"testing"

	"time"

	"github.com/kr/pretty"
)

func TestParallel(t *testing.T) {
	Parallel([]func(func(error, ...interface{})){
		func(cb func(error, ...interface{})) {
			<-time.After(200 * time.Millisecond)
			cb(nil, "aa")
		},
		func(cb func(error, ...interface{})) {
			<-time.After(100 * time.Millisecond)
			cb(nil, "bb")
		},
		// func(cb func(error, ...interface{})) {
		// 	<-time.After(100 * time.Millisecond)
		// 	cb(errors.New("dog"), nil)
		// },
	}, func(res []interface{}, err error) {
		if err != nil {
			pretty.Println(err)
			return
		}

		pretty.Println(res)
	})

	<-time.After(2 * time.Second)
}
