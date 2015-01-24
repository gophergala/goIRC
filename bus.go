package main

import (
	"fmt"
)

type EventType int

const (
	EventSay EventType = iota
	EventLeave
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
	switch event.event_type {
	case EventSay:
		fmt.Printf("%q(%d)> %q\n", s.name, event.event_type, event.event_data)
	case EventLeave:
		fmt.Printf("%q(%d)> LEAVING!!!\n", s.name, event.event_type)
	}
}

func (bus *EventBus) Publish(event *Event) {
	fmt.Printf("\npublishing -%d- data: %q\n", event.event_type, event.event_data)
	for _, subscriber := range bus.subscribers {
		subscriber.OnEvent(event) //currently slower than without the goroutine
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

	for i := 0; i < 10; i++ {
		s := Subscriber{name: fmt.Sprintf("client_%v", i)}
		bus.Subscribe(&s)
	}

	bus.Publish(&e)
	bus.Publish(&Event{event_type: EventLeave})
}
