package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type ConnectionStatus int

const (
	SocketConnected ConnectionStatus = iota
	UserPassSent
	UserNickSent
	UserUserInfoSent
	UserRegistered
)

type User struct {
	Nick     string
	Ident    string
	RealName string
	Conn     net.Conn
	Status   ConnectionStatus
}

func handleConnection(conn net.Conn, buses map[string]*EventBus) {
	client := User{Status: UserPassSent, Conn: conn}
	reader := bufio.NewReader(conn)

	for {
		status, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		status = strings.TrimSpace(status)

		if client.Status < UserRegistered {
			regCmd := strings.SplitN(status, " ", 2)

			switch regCmd[0] {
			case "NICK":
				client.Nick = regCmd[1]
				conn.Write([]byte("welcome " + client.Nick + "\n"))
				client.Status = UserNickSent
			case "USER":
				var uname, hname, sname, rname string
				fmt.Sscanf(regCmd[1], "%q %q %q :%q", uname, hname, sname, rname) //TODO(jz) need to split on : in case real name has spaces
				fmt.Println(hname + uname)                                        //just so we don't get the unused var error
				client.Ident = uname
				client.RealName = rname
				client.Status = UserRegistered
			case "PASS":
				client.Nick = regCmd[1]
				client.Ident = regCmd[1]
				client.RealName = regCmd[1]
				client.Status = UserRegistered

			default:
				conn.Write([]byte("you must register first. try nick or user?\n"))
			}

		} else {
			// split <command> <target>:<data>
			var cmd, target, data string
			s := strings.SplitN(status, ":", 2)
			_, err = fmt.Sscanf(s[0], "%s %s", &cmd, &target)
			if err != nil {
				fmt.Println(err)
				conn.Write([]byte("Invalid input! CHECK YOUR(self) SYNTAX\n"))
				continue
			}
			fmt.Println(len(s))
			if len(s) == 2 {
				data = s[1]
			} else {
				conn.Write([]byte("SYNTAX...PLEASE....\n"))
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
				b.Subscribe(UserJoin, &client)
				b.Subscribe(PrivMsg, &client)

				message := fmt.Sprintf("%s joined %s!\n", client.Nick, target)
				b.Publish(&Event{UserJoin, message})
			case "MSG":
				b, ok := buses[target]
				if !ok {
					conn.Write([]byte("Channel does not exist\n"))
				}
				// implment check if client is subscribed to channel here
				message := fmt.Sprintf("%s: %s", client.Nick, data)
				b.Publish(&Event{PrivMsg, message})
			default:
				conn.Write([]byte("No Command match\n"))
			}
		}
	}
}
