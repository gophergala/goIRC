package main

import (
	"fmt"
)

type EventType int

const (
	ChannelUserJoin EventType = iota
	ChannelUserPart
	ChannelPrivMsg
)

type EventBus struct {
	subscribers map[EventType][]*Subscriber
	channel     *Channel
}

type Event struct {
	event_type EventType
	event_data string
}

type Subscriber struct {
	name string
}

type Channel struct {
	name  string
	topic string
}

func (s *Subscriber) OnEvent(event *Event) {
	switch event.event_type {
	case ChannelUserJoin:
		fmt.Printf("%q(%d)> %q\n", s.name, event.event_type, event.event_data)
	}
}

func (bus *EventBus) Publish(event *Event) {
	fmt.Printf("\npublishing -%d- data: %q\n", event.event_type, event.event_data)
	for _, subscriber := range bus.subscribers[event.event_type] {
		subscriber.OnEvent(event) //currently slower than without the goroutine
	}
	fmt.Println("done publishing")
}

func (bus *EventBus) Subscribe(event_type EventType, subscriber *Subscriber) {
	bus.subscribers[event_type] = append(bus.subscribers[event_type], subscriber)
}

func main() {

	gophers := Channel{name: "#gophers", topic: "gogo gophergala!"}

	buses := make(map[string]EventBus)
	buses["#gophers"] = EventBus{channel: &gophers}
	fmt.Println("name: " + buses["#gophers"].channel.name) // bus := EventBus{}
	sub := Subscriber{name: "a_client"}
	var b EventBus //go can't infer buses["key"] is an EventBus?
	b = buses["#gophers"]
	b.Subscribe(ChannelUserJoin, &sub)

	// e := Event{event_type: EventSay, event_data: "hello, world"}

	for i := 0; i < 10; i++ {
		s := Subscriber{name: fmt.Sprintf("client_%v", i)}
		b.Subscribe(ChannelUserJoin, &s)
	}
	e := Event{event_type: ChannelUserJoin, event_data: "Alvin has joined!"}
	b.Publish(&e)
	// bus.Publish(&Event{event_type: EventLeave})
}
