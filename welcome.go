package main

/*
 * welcome.go
 * Welcome the user to the channel
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160821
 */

import (
	"fmt"
	"log"
)

const (
	/* WELCOMETEXT is the text to send a user the first time he's seen */
	WELCOMETEXT = `Welcome, %v.  This is a family-friendly channel.`
)

/* WelcomeUser welcomes a user the first time he's seen */
func WelcomeUser(nick, channel string) error {
	/* Welcome the user */
	if err := Privmsg(channel, fmt.Sprintf(WELCOMETEXT, nick)); nil != err {
		return err
	}
	log.Printf("Welcomed %v", nick)
	return nil
}
