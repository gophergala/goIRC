package main

import (
	"fmt"
)

type EventType int

const (
	EventSay   EventType = iota
	EventLeave EventType
)

type EventBus struct {
	subscribers []*Subscriber
}

type Event struct {
	event_type EventType
	event_data string
}

type Subscriber struct {
	name string
}

func (s *Subscriber) OnEvent(event *Event) {
	fmt.Printf("%q(%q)> %q\n", s.name, event.event_type, event.event_data)
}

func (bus *EventBus) Publish(event *Event) {
	fmt.Printf("\npublishing -%q- data: %q\n", event.event_type, event.event_data)
	for _, subscriber := range bus.subscribers {
		go subscriber.OnEvent(event) //currently slower than without the goroutine
	}
	fmt.Println("done publishing")
}

func (bus *EventBus) Subscribe(subscriber *Subscriber) {
	bus.subscribers = append(bus.subscribers, subscriber)
}

func main() {
	bus := EventBus{}

	sub := Subscriber{name: "a_client"}
	bus.Subscribe(&sub)

	e := Event{event_type: EventSay, event_data: "hello, world"}

	for i := 0; i < 1000000; i++ {
		sub = Subscriber{name: fmt.Sprintf("client_%v", i)}
		bus.Subscribe(&sub)
	}

	bus.Publish(&e)
}
