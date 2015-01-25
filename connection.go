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
		fmt.Println("READY TO READ!!!!!!!!!!!!!!!!!!!!!!")

		status, err := reader.ReadString('\n')
		fmt.Println("GOT BY READ!!!!!!!!!!!!!!!!!!!!!!")
		if err != nil {
			panic(err)
		}

		if client.Status < UserRegistered {
			regCmd := strings.Split(status, " ")

			switch regCmd[0] {
			case "NICK":
				client.Nick = regCmd[1]
				conn.Write([]byte("welcome " + client.Nick + "\r\n"))
				client.Status = UserNickSent

			//}
			default:
				conn.Write([]byte("you must register first. try nick?\r\n"))
			}

		} else {
			var cmd, target, data string

			n, err := fmt.Sscanf(status, "%s %s %q", &cmd, &target, &data)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(n)
			fmt.Println(cmd, target, data)

			switch cmd {
			case "JOIN":
				b := buses[target]
				b.Subscribe(ChannelUserJoin, &client)
				b.Subscribe(ChannelMsg, &client)
				message := fmt.Sprintf("%s joined %s!\n", client.Nick, target)
				b.Publish(&Event{ChannelUserJoin, message})
				fmt.Println("JOIN CASE WAS HIT!!!!!!!!!!!!")
			case "MSG":
				b := buses[target]
				message := fmt.Sprintf("%s: %s\n", client.Nick, data)
				b.Publish(&Event{ChannelMsg, message})
				fmt.Println("MSG WAS HIT!!!!!")
			}
		}
	}
}
