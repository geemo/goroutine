package race

import (
	"test/goroutine/race/channel"
	"testing"
	"time"

	"github.com/kr/pretty"
)

func TestRace(t *testing.T) {
	vals, err := Race(
		2,             // 返回的参数个数
		2*time.Second, // 返回的deadline
		func(ch channel.Channel) {
			<-time.After(200 * time.Millisecond)
			ch.Send("hello", "geemo")
			ch.Done()
		},
		func(ch channel.Channel) {
			<-time.After(100 * time.Millisecond)
			// ch.Send("geemo", "hello")
			ch.Send("geemo")
			<-time.After(50 * time.Millisecond)
			ch.Send("hello")
			ch.Done()
		},
		func(ch channel.Channel) {
			<-time.After(160 * time.Millisecond)
			ch.Send("hello", "world")
			ch.Done()
		},
	)

	pretty.Println(vals, err)
}
