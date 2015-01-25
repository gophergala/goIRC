package main

import (
	"net"
)

type ConnectionStatus int

const (
	SocketConnected ConnectionStatus = iota
	UserPassSent
	UserNickSent
	UserUserInfoSent
	UserRegistered
)

func handleConnection(conn net.Conn, buses *MasterBus) {
	client := User{ConnectionStatus: SocketConnected}

	for {
		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			panic("OH NOEESssss")
		}
		if client.Status < UserRegistered {
			switch cmd {
			case "NICK":
				client.Nick = strings.Split(nick, " ")[1]
				if client.Status >= UserPassSent {
					client.Status = UserNickSent
				}
			case "USER":

			}
		} else {

			var cmd, target, data string
			n, err := fmt.Sscanf(status, "%s %s %q", &cmd, &target, &data)
			fmt.Println(n)
			fmt.Println(cmd, target, data)

			switch cmd {
			case "JOIN":
				client = User{Nick: data, conn: conn}
				b := buses[target]
				b.Subscribe(ChannelUserJoin, &client)
				b.Subscribe(ChannelMsg, &client)
				message := fmt.Sprintf("%s joined %s!\n", client.Nick, target)
				b.Publish(&Event{ChannelUserJoin, message})
			case "MSG":
				b := buses[target]
				message := fmt.Sprintf("%s: %s\n", client.Nick, data)
				b.Publish(&Event{ChannelMsg, message})
			}
		}
		// this just echos whatever is sent over
		//n, err := conn.Write([]byte(status))
		// if err != nil {
		// 	fmt.Println("Not looking too good")
		// }
		// fmt.Println(n)
	}
}
