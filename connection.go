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

	commands := make(map[string]func(map[string]*EventBus, *User, string, string))
	commands["JOIN"] = handleJoin
	commands["PRIVMSG"] = handleMsg
	commands["NICK"] = handleNick
	commands["PART"] = handlePart
	commands["HELP"] = handleHelp

	for {
		status, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		status = strings.TrimSpace(status)

		// allows user to enter empty strings
		if len(status) == 0 {
			conn.Write([]byte(""))
			continue
		}

		if client.Status < UserRegistered {
			regCmd := strings.SplitN(status, " ", 2)

			switch regCmd[0] {
			case "NICK":
				//client.Nick = regCmd[1]
				handleNick(buses, &client, regCmd[1], "")
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
			// split <command> <target> :<data>

			var cmd, target, data string
			s := strings.SplitN(status, ":", 2)
			if len(s) > 1 {
				data = s[1]
			}
			_, err = fmt.Sscanf(s[0], "%s %s", &cmd, &target)
			if err != nil {
				fmt.Println(err)
				conn.Write([]byte("Invalid input! CHECK YOUR(self) SYNTAX\n"))
				continue
			}

			if _, ok := commands[cmd]; ok {
				commands[cmd](buses, &client, target, data)
			}
		}
	}
}

func handlePart(buses map[string]*EventBus, client *User, target string, data string) {
	message := fmt.Sprintf("%s parted %s!\n", client.Nick, target)
	buses[target].Publish(&Event{event_type: UserPart, event_data: message})
	buses[target].Unsubscribe(UserPart, client)
	buses[target].Unsubscribe(UserJoin, client)
	buses[target].Unsubscribe(PrivMsg, client)
}

func handleJoin(buses map[string]*EventBus, client *User, target string, data string) {
	fmt.Println("!!!!!!!!! JOIN")
	b, ok := buses[target]
	if !ok {
		// need to add support for channel topic
		newChannel := Channel{name: target, topic: "gogo new channel!"}
		buses[newChannel.name] = &EventBus{subscribers: make(map[EventType][]Subscriber), channel: &newChannel}
		b = buses[newChannel.name]
	}
	b.Subscribe(UserJoin, client)
	b.Subscribe(PrivMsg, client)
	b.Subscribe(UserPart, client)

	message := fmt.Sprintf("%s joined %s!\n", client.Nick, target)
	b.Publish(&Event{UserJoin, message})
}

func handleNick(buses map[string]*EventBus, client *User, target string, data string) {
	client.Nick = target
	client.Conn.Write([]byte("nick set to:" + client.Nick + "\n"))
}

func handleMsg(buses map[string]*EventBus, client *User, target string, data string) {
	b, ok := buses[target]
	if !ok {
		client.Conn.Write([]byte("Channel does not exist\n"))
	}
	// implment check if client is subscribed to channel here
	message := fmt.Sprintf("(%s)%s: %s\n", target, client.Nick, data)
	b.Publish(&Event{event_type: PrivMsg, event_data: message})

}

func handleHelp(buses map[string]*EventBus, client *User, target string, data string) {
	k, ok := Help[target]
	if !ok {
		client.Conn.Write([]byte("Available Commands:\n"))
		for h := range Help {
			client.Conn.Write([]byte(h + "\n"))
		}
	} else {
		client.Conn.Write([]byte("Summary: " + k.Summary + "\nUsage: " + k.Syntax + "\n"))
	}
}
