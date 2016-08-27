package main

/*
 * privmsg.go
 * Handle private messages
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160826
 */

import (
	"fmt"
	"log"
	"strings"
)

/* HandlePrivmsg is called to handle a private message */
func HandlePrivmsg(nick, op, msg string) error {
	go handlePrivmsg(nick, op, msg)
	return nil
}

/* handlePrivmsg does all the real work for HandlePrivmsg */
func handlePrivmsg(nick, op, msg string) {
	/* Work out real sender */
	parts := strings.SplitN(msg, " ", 2)
	if 2 != len(parts) {
		log.Printf("Short privmsg from %v: %q", nick, msg)
	}
	tgt := parts[0]
	msg = parts[1][1:]

	/* Work out to whom to send the reply */
	var replyto string
	if strings.HasPrefix(tgt, "#") {
		replyto = tgt
		go WelcomeUser(nick, tgt)
	} else {
		replyto = nick
	}

	/* Log the first and latest privmsg time */
	go LogFirstLast(
		nick,
		"PRIVMSG",
		fmt.Sprintf("%v %v", replyto, msg),
	)

	log.Printf("[PRIVMSG] %v (reply-to %v): %v", nick, replyto, msg)

	/* Handle commands */
	/* TODO: Un-hardcode this */
	if strings.HasPrefix(msg, "!") {
		if err := HandleCommand(msg[1:], nick, replyto); nil != err {
			log.Printf(
				"Error handling command %q from %v in %q: %v",
				msg,
				nick,
				replyto,
				err,
			)
		}
	}
}
