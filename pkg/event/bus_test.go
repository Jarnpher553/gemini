package event

import (
	"github.com/Jarnpher553/gemini/pkg/redis"
	"testing"
)

func TestSubscribe(t *testing.T) {
	rd := redis.New(redis.Addr("192.168.95.139:6379"), redis.DB(15), redis.PoolSize(100), redis.Pwd("shinssonshinsson1234"))
	Bind(rd)

	_ = Subscribe("event/demo")

	go func() {
		Publish("event/demo", NewEvent("do", 3))

		Publish("event/demo", NewEvent("do", 3))

		Publish("event/demo", NewEvent("do", 3))

		close(bus.ch)
	}()

	for ev := range Events() {
		t.Log(ev)
	}

	t.Log("end")
}
