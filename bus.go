package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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
	conn net.Conn
}

type Channel struct {
	name  string
	topic string
}

// something funky going on here
// type Subscriber interface {
// 	OnEvent(event *Event)
// }

func (s *Subscriber) OnEvent(event *Event) {
	switch event.event_type {
	case ChannelUserJoin:
		//fmt.Printf("%q(%d)> %q\n", s.Nick, event.event_type, event.event_data)
		_, err := s.conn.Write([]byte(event.event_data))
		if err != nil {
			fmt.Println("Not looking too good")
		}
	case ChannelMsg:
		_, err := s.conn.Write([]byte(event.event_data))
		if err != nil {
			fmt.Println("Not looking too good")
		}
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

func handleConnection(conn net.Conn, buses map[string]*EventBus) {
	var client Subscriber
	for {
		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			panic("OH NOEESssss")
		}
		var cmd, target, data string

		// split <command> <target>:<data>
		s := strings.SplitN(status, ":", 2)
		_, err = fmt.Sscanf(s[0], "%s %s", &cmd, &target)
		if err != nil {
			panic(err)
		}
		data = s[1]

		switch cmd {
		case "JOIN":
			b, ok := buses[target]
			if !ok {
				// need to add support for channel topic
				newChannel := Channel{name: target, topic: "gogo new channel!"}
				buses[newChannel.name] = &EventBus{make(map[EventType][]*Subscriber), &newChannel}
				b = buses[newChannel.name]
			}
			data = data[:len(data)-2]
			fmt.Println(data)
			client = Subscriber{Nick: data, conn: conn}
			b.Subscribe(ChannelUserJoin, &client)
			b.Subscribe(ChannelMsg, &client)

			message := fmt.Sprintf("%s joined %s!\n", client.Nick, target)
			b.Publish(&Event{ChannelUserJoin, message})
		case "MSG":
			b := buses[target]
			message := fmt.Sprintf("%s: %s", client.Nick, data)
			b.Publish(&Event{ChannelMsg, message})
		}
	}
}

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
