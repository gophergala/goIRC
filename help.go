package main

var Help = map[string]string {
	CMD_JOIN: 	"join a new channel /n
				syntax: JOIN <#channel>: <username>"
	CMD_MSG: 	"send a message to a channel /n
				syntax: <MSG <#channel>: <your message here>"
	
	// not yet implemented
	CMD_TOPIC: 	"see the topic for a certain channel /n
				syntax: TOPIC <#channel>: <your topic here>"
	CMD_USERS: 	"get a list of all users in a channel /n
				syntax: USERS <#channel>"
	CMD_NICK: 	"change your nickname /n
				syntax: NICK <new_nick_here>"
	CMD_QUIT:	"quit the server
				syntax: <QUIT>"
	CMD_PART:	"leave the channel
				syntax: PART <#channel>"
	CMD_KICK:	"kick out a user from the channel - ops only /n
				syntax: KICK <#channel>: <user_to_kick>"
	CMD_INVITE:	"invite a user to the channel /n
				syntax: INVITE <#channel>: <user_to_invite>"
	CMD_KILL:	"disconnect a user - ops only /n
				syntax: KILL <user_to_kill>"
}

