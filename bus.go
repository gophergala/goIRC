package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	// _ "sync"
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

type User struct {
	Nick     string
	Ident    string
	RealName string
	Conn     net.Conn
	// Status   ConnectionStatus
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

func handleConnection(Conn net.Conn, buses map[string]*EventBus) {
	var client User
	for {
		status, err := bufio.NewReader(Conn).ReadString('\n')
		if err != nil {
			panic("OH NOEESssss")
		}
		var cmd, target, data string

		// split <command> <target>:<data>
		s := strings.SplitN(status, ":", 2)
		_, err = fmt.Sscanf(s[0], "%s %s", &cmd, &target)
		if err != nil {
			fmt.Println(err)
			Conn.Write([]byte("Invalid input! CHECK YOUR(self) SYNTAX\n"))
			continue
		}
		fmt.Println(len(s))
		if len(s) == 2 {
			data = s[1]
		} else {
			Conn.Write([]byte("SYNTAX...PLEASE....\n"))
			continue
		}

		switch cmd {
		case "JOIN":
			b, ok := buses[target]
			if !ok {
				// need to add support for channel topic
				newChannel := Channel{name: target, topic: "gogo new channel!"}
				buses[newChannel.name] = &EventBus{make(map[EventType][]Subscriber), &newChannel}
				b = buses[newChannel.name]
			}
			data = data[:len(data)-2]
			fmt.Println(data)
			client = User{Nick: data, Conn: Conn}
			b.Subscribe(UserJoin, &client)
			b.Subscribe(PrivMsg, &client)

			message := fmt.Sprintf("%s joined %s!\n", client.Nick, target)
			b.Publish(&Event{UserJoin, message})
		case "MSG":
			b, ok := buses[target]
			if !ok {
				Conn.Write([]byte("Channel does not exist\n"))
			}
			// implment check if client is subscribed to channel here
			message := fmt.Sprintf("%s: %s", client.Nick, data)
			b.Publish(&Event{PrivMsg, message})
		default:
			Conn.Write([]byte("No Command match\n"))
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
