package main

import (
	"bufio"
	"fmt"
	"net"
)

type EventType int

const (
	ChannelUserJoin EventType = iota
	ChannelUserPart
	ChannelPrivMsg
	ChannelMsg
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
	Nick string
}

type Channel struct {
	name  string
	topic string
}

// something funky going on here
// type Subscriber interface {
// 	OnEvent(event *Event)
// }

func (s *Subscriber) OnEvent(event *Event, conn net.Conn) {
	switch event.event_type {
	case ChannelUserJoin:
		//fmt.Printf("%q(%d)> %q\n", s.Nick, event.event_type, event.event_data)
		_, err := conn.Write([]byte(s.Nick + " joined " + event.event_data))
		if err != nil {
			fmt.Println("Not looking too good")
		}
	case ChannelMsg:
		_, err := conn.Write([]byte(event.event_data))
		if err != nil {
			fmt.Println("Not looking too good")
		}
	}
}
func (bus *EventBus) Publish(event *Event, conn net.Conn) {
	fmt.Printf("\npublishing -%d- data: %q\n", event.event_type, event.event_data)
	for _, subscriber := range bus.subscribers[event.event_type] {
		subscriber.OnEvent(event, conn) //currently slower than without the goroutine
	}
	fmt.Println("done publishing")
}

func (bus *EventBus) Subscribe(event_type EventType, subscriber *Subscriber) {
	bus.subscribers[event_type] = append(bus.subscribers[event_type], subscriber)
}

func handleConnection(conn net.Conn) {
	for {
		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			panic("OH NOEESssss")
		}
		var cmd, target, data string
		n, err := fmt.Sscanf(status, "%s %s %q", &cmd, &target, &data)
		fmt.Println(n)
		fmt.Println(cmd, target, data)

		// this does not realy work...
		switch cmd {
		case "JOIN":
			b := buses[target]
			b.Subscribe(ChannelUserJoin, &Subscriber{Nick: data})
			b.Publish(&Event{ChannelUserJoin, target}, conn)
		case "MSG":
			b := buses[target]
			b.Publish(&Event{ChannelMsg, data}, conn)
		}
		// this just echos whatever is sent over
		//n, err := conn.Write([]byte(status))
		// if err != nil {
		// 	fmt.Println("Not looking too good")
		// }
		// fmt.Println(n)
	}
}

var buses map[string]*EventBus

func init() {
	// init event bus map
	buses = make(map[string]*EventBus)

	// make new channel #gophers
	gophers := Channel{name: "#gophers", topic: "gogo gophergala!"}

	buses[gophers.name] = &EventBus{make(map[EventType][]*Subscriber), &gophers}
	fmt.Println("New Channel: " + buses[gophers.name].channel.name)
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
	ln, err := net.Listen("tcp", ":3030")
	if err != nil {
		panic("Listen not WORKING")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic("nope not Accepting")
		}
		go handleConnection(conn)
	}

}
