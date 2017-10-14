package main

/*
 * part.go
 * Greet users when they join the channel
 * By MagisterQuis
 * Created 20170913
 * Last Modified 20170913
 */

import (
	"fmt"
	"io"
	"log"
	"strings"
)

/* handlePart tells the channel that the user left */
func handlePart(w io.Writer, channel, nick string) error {
	LOCK.Lock()
	defer LOCK.Unlock()
	/* Don't say goodbye if we're a guest */
	if _, ok := GUESTCHANNELS[channel]; ok {
		return nil
	}
	_, err := fmt.Fprintf(w, "PRIVMSG %v :%q left\r\n", channel, nick)
	log.Printf("[GOODBYE] %v (%v)", nick, channel)
	return err
	/* TODO: Finish this */
}

// Part leaves a channel
func Part(w io.Writer, channel string) {
	/* Make sure the channel is a channel */
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}
	/* Ask to part the channel */
	if _, err := fmt.Fprintf(w, "PART %v\r\n", channel); nil != err {
		log.Printf("[PART] Error parting %q: %v", channel, err)
		return
	}
	log.Printf("[PART] %q", channel)
}
