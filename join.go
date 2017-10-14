package main

/*
 * join.go
 * Greet users when they join the channel
 * By MagisterQuis
 * Created 20170625
 * Last Modified 20170913
 */

import (
	"fmt"
	"io"
	"log"
	"strings"
)

// DEFAULTGREETING is the greeting sent to those without a special greeting.
const DEFAULTGREETING = "Welcome!  This is a family-friendly channel."

/* handleJoin greets users who join the channel.  It tries to use a greeting
from GREETINGS, and failing that,  it sends the nick, a colon, and the default
greeting. */
func handleJoin(w io.Writer, channel, nick string) error {
	nick = strings.ToLower(nick)

	/* Get the nick-apprpriate greeting */
	LOCK.Lock()
	_, isBot := BOTS[nick]
	_, guest := GUESTCHANNELS[channel]
	greeting, ok := GREETINGS[strings.ToLower(nick)]
	LOCK.Unlock()

	/* No greeting if we're a guest */

	/* No greeting for bots */
	if isBot || guest {
		return nil
	}

	/* Make one if none exists */
	if !ok {
		greeting = fmt.Sprintf("%v: %v", nick, DEFAULTGREETING)
	}

	/* Try and send it */
	_, err := fmt.Fprintf(w, "PRIVMSG %v :%v\r\n", channel, greeting)
	log.Printf("[GREETING] %v (%v)", nick, greeting)
	return err
}

// Join joins a channel
func Join(w io.Writer, channel string) {
	/* Make sure the channel is a channel */
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}
	/* Try to join the channel */
	if _, err := fmt.Fprintf(w, "JOIN %v\r\n", channel); nil != err {
		log.Printf("[JOIN] Error joining %q: %v", channel, err)
		return
	}
	log.Printf("[JOIN] %q", channel)
}
