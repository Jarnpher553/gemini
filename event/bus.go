package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jarnpher553/gemini/redis"
	"github.com/Jarnpher553/gemini/shortuuid/snow"
	"time"
)

type Bus struct {
	*redis.RdClient
	ch   chan Event
	name string
}

type Event struct {
	ID        string
	Timestamp int64
	Action    string
	Content   interface{}
}

var bus = &Bus{}

func NewEvent(action string, content interface{}) *Event {
	return &Event{
		ID:        fmt.Sprintf("ev_%s", snow.NextID()),
		Timestamp: time.Now().UnixNano() / 1e6,
		Action:    action,
		Content:   content,
	}
}

func Bind(client *redis.RdClient) {
	bus.RdClient = client
	bus.ch = make(chan Event, 100)
}

func Subscribe(name string) error {
	if bus.name != "" {
		return errors.New("event bus has existed")
	}
	bus.name = name
	ps := bus.RdClient.Subscribe(name)
	go func() {
		for message := range ps.Channel() {
			var ev Event
			_ = json.Unmarshal([]byte(message.Payload), &ev)
			bus.ch <- ev
		}
	}()
	return nil
}

func Events() <-chan Event {
	return bus.ch
}

func Publish(name string, ev *Event) error {
	marshal, _ := json.Marshal(ev)
	ret := bus.RdClient.Publish(name, marshal)
	if !ret {
		return errors.New("publish error")
	}
	return nil
}
