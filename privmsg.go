package main

/*
 * privmsg.go
 * Handle triggered sayings
 * by MagisterQuis
 * Created 20170625
 * Last Modified 20170625
 */

import (
	"fmt"
	"io"
	"log"
	"strings"
)

/* handlePrivmsg responds to a privmsg from src if there's a canned
response */
func handlePrivmsg(w io.Writer, nick, msg string) error {
	/* Break apart message */
	parts := strings.SplitN(msg, " :", 2)
	if 2 != len(parts) {
		return fmt.Errorf("not enough parts to %q", msg)
	}
	var (
		channel = strings.TrimSpace(parts[0])
		trigger = strings.TrimSpace(parts[1])
	)

	/* See if we have a response */
	res, found := getResponse(trigger)
	if !found {
		return nil
	}

	/* If so, send it back */
	if _, err := fmt.Fprintf(
		w,
		"PRIVMSG %v :%v: %v\r\n",
		channel,
		nick,
		res,
	); nil != err {
		return err
	}
	log.Printf("[COMMAND] %v: %v -> %v", nick, trigger, res)

	return nil
}

/* getResponse returns the response for the trigger t, as well as whether there
actually was one */
func getResponse(t string) (string, bool) {
	LOCK.Lock()
	defer LOCK.Unlock()
	r, ok := RESPONSES[t]
	return r, ok
}
