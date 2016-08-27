package main

/*
 * welcome.go
 * Welcome the user to the channel
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160826
 */

import "log"

const (
	/* WELCOMETEXT is the text to send a user the first time he's seen */
	WELCOMETEXT = `Welcome, %v.  This is a family-friendly channel.`
)

/* WelcomeUser welcomes a user the first time he's seen */
func WelcomeUser(nick, channel string) {
	/* Don't bother with bots */
	if _, ok := BOTS[nick]; ok {
		return
	}
	/* Note we've seen user */
	seen, err := PutBool(sb("welcomed"), true, sb("viewers"), sb(nick))
	if nil != err {
		log.Printf(
			"Unable to note welcome of %v in %v",
			nick,
			channel,
			err,
		)
	}
	/* If we'd already seen him, return */
	if nil != seen && *seen {
		return
	}
	/* If not, welcome the user */
	if err := Privmsg(channel, WELCOMETEXT, nick); nil != err {
		log.Printf(
			"Unable to welcome %v to %v: %v",
			nick,
			channel,
			err,
		)
	}
	log.Printf("[WELCOME] %v (%v)", nick, channel)
	return
}

/* TODO: Put in appropriate places */
