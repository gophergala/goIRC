package main

type HelpCommand struct {
	Summary string
	Syntax  string
	// Help		string - more detailed help instructions
	// Ops		boolean - flag to see who can use the command
}

var Help = map[string]HelpCommand{
	"CMD_JOIN": HelpCommand{
		Summary: "join a new channel",
		Syntax:  "JOIN <#channel>: <username>",
	},
	"CMD_MSG": HelpCommand{
		Summary: "send a message to a channel",
		Syntax:  "MSG <#channel>: <your message here>",
	},

	// not yet implemented
	"CMD_TOPIC": HelpCommand{
		Summary: "see the topic for a certain channel",
		Syntax:  "TOPIC <#channel>: <your topic here>",
	},
	"CMD_USERS": HelpCommand{
		Summary: "get a list of all users in a channel",
		Syntax:  "syntax: USERS <#channel>",
	},
	"CMD_NICK": HelpCommand{
		Summary: "change your nickname",
		Syntax:  "NICK <new_nick_here>",
	},
	"CMD_QUIT": HelpCommand{
		Summary: "quit the server",
		Syntax:  "<QUIT>",
	},
	"CMD_PART": HelpCommand{
		Summary: "leave the channel",
		Syntax:  "PART <#channel>",
	},
	"CMD_KICK": HelpCommand{
		Summary: "kick out a user from the channel - ops only",
		Syntax:  "KICK <#channel>: <user_to_kick>",
	},
	"CMD_INVITE": HelpCommand{
		Summary: "invite a user to the channel",
		Syntax:  "INVITE <#channel>: <user_to_invite>",
	},
	"CMD_KILL": HelpCommand{
		Summary: "disconnect a user - ops only",
		Syntax:  "KILL <user_to_kill>",
	},
}
