package timeout

import (
	"testing"
	"time"

	"github.com/kr/pretty"
)

func TestSetTimeout(t *testing.T) {
	SetTimeout(func() {
		pretty.Println("hello geemo")
	}, 1*time.Second)

	<-time.After(2 * time.Second)
}
