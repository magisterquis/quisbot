package main

import "log"

/*
 * part.go
 * Handle PART messages
 * By J. Stuart McMurray
 * Created 20160826
 * Last Modified 20160826
 */

/* handlePart handles PART messages */
func HandlePart(nick, op, channel string) error {
	/* Log the first and latest PART */
	go LogFirstLast(nick, op, channel)
	log.Printf("[PART] %v from %v", nick, channel)
	return nil
}
