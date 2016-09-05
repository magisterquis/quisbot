package main

/*
 * commands.go
 * User-facing commands
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160904
 */

import (
	"fmt"
	"strings"
)

const IGN = "magisterquis"

/* Commands which return simple information */
var simpleCommands = map[string]string{
	"ign": fmt.Sprintf("His IGN is %v", IGN),
}

var commandFunctions = map[string]func(nick, replyto, args string) error{
	"bet":     PlaceBet,
	"balance": CheckBalance,
	/* TODO: !sh */
}

/* HandleCommand handles a command starting with a ! */
func HandleCommand(msg, nick, replyto string) error {
	if "" == msg {
		return nil
	}
	/* Split into command and argument */
	parts := strings.SplitN(msg, " ", 2)
	cmd := strings.TrimSpace(parts[0])
	var args string
	if 2 == len(parts) {
		args = strings.TrimSpace(parts[1])
	}

	/* TODO: Help message */

	/* Check the simple ones */
	if info, ok := simpleCommands[cmd]; ok {
		go Privmsg(replyto, fmt.Sprintf("%v: %v", nick, info))
		return nil
	}

	/* Check the functions */
	if f, ok := commandFunctions[cmd]; ok {
		return f(nick, replyto, args)
	}

	return nil
}
