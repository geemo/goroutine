package interval

import (
	"testing"
	"time"

	"github.com/kr/pretty"
)

func TestSetInterval(t *testing.T) {
	SetInterval(func() {
		pretty.Println("hello geemo: ", time.Now())
	}, 1*time.Second)

	<-time.After(5 * time.Second)
}
