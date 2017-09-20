package eachLimit

import (
	"errors"
	"testing"

	"time"

	"github.com/kr/pretty"
)

func TestEachLimit(t *testing.T) {
	arr := []interface{}{1, 2, 3, 4, 5}
	EachLimit(
		arr,
		2,
		func(elem interface{}, next func(error)) {
			pretty.Println(elem)

			val, ok := elem.(int)
			if !ok {
				next(errors.New("conv error"))
				return
			}

			if val == 2 {
				next(errors.New("stop"))
				return
			}

			next(nil)
		},
		func(err error) {
			pretty.Println("done")

			if err != nil {
				pretty.Println(err)
				return
			}
		},
	)

	<-time.After(1 * time.Second)
}
