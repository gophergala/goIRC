package main

import (
	"fmt"
	"net"
	_ "sync"
)

type EventType int

const (
	UserJoin EventType = iota
	UserPart
	PrivMsg
)

type Subscriber interface {
	OnEvent(*Event)
}

type EventBus struct {
	subscribers map[EventType][]Subscriber
	channel     *Channel
}

type Event struct {
	event_type EventType
	event_data string
}

type Channel struct {
	name  string
	topic string
}

func (u *User) OnEvent(event *Event) {
	switch event.event_type {
	case UserJoin:
		//fmt.Printf("%q(%d)> %q\n", s.Nick, event.event_type, event.event_data)
		_, err := u.Conn.Write([]byte(event.event_data))
		if err != nil {
			fmt.Println("Not looking too good")
		}
	case PrivMsg:
		_, err := u.Conn.Write([]byte(event.event_data))
		if err != nil {
			fmt.Println("Not looking too good")
		}
	}
}
func (bus *EventBus) Publish(event *Event) {
	fmt.Printf("\npublishing -%d- data: %q\n", event.event_type, event.event_data)
	for _, subscriber := range bus.subscribers[event.event_type] {
		go subscriber.OnEvent(event) //currently slower than without the goroutine
	}
	fmt.Println("done publishing")
}

func (bus *EventBus) Subscribe(event_type EventType, subscriber Subscriber) {
	bus.subscribers[event_type] = append(bus.subscribers[event_type], subscriber)
}

var buses map[string]*EventBus

func init() {

	// make new channel #gophers
	// gophers := Channel{name: "#gophers", topic: "gogo gophergala!"}

	// buses[gophers.name] = &EventBus{make(map[EventType][]*Subscriber), &gophers}
	// fmt.Println("New Channel: " + buses[gophers.name].channel.name)
	// sub := Subscriber{Nick: "a_client"}
	// fmt.Println("New Subscriber: " + sub.Nick)

	// b := buses["#gophers"]
	// b.Subscribe(ChannelUserJoin, &sub)

	// // e := Event{event_type: EventSay, event_data: "hello, world"}

	// for i := 0; i < 10; i++ {
	// 	s := Subscriber{Nick: fmt.Sprintf("client_%v", i)}
	// 	b.Subscribe(ChannelUserJoin, &s)
	// }
	// e := Event{event_type: ChannelUserJoin, event_data: "Alvin has joined!"}
	//b.Publish(&e)
	// bus.Publish(&Event{event_type: EventLeave})
}
func main() {
	// init event bus map
	buses := make(map[string]*EventBus)

	ln, err := net.Listen("tcp", ":3030")

	if err != nil {
		panic("Listen not WORKING")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic("nope not Accepting")
		}
		go handleConnection(conn, buses)
	}

}
