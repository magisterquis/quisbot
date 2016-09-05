package main

import "log"

/*
 * join.go
 * Handle JOIN messages
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160826
 */

/* HailList is a list of users to hail on JOIN */
var HailList = []string{
	"mairielpis",
	"imperatorpat",
	"h46u",
}

/* hails is the set of users to hail on JOIN */
var hails map[string]struct{}

/* init makes hails from HailList */
func init() {
	hails = make(map[string]struct{})
	for _, n := range HailList {
		hails[n] = struct{}{}
	}
}

/* handleJoin handles JOIN messages */
func HandleJoin(nick, op, channel string) error {
	/* Log the first and latest join */
	go LogFirstLast(nick, op, channel)
	/* Welcome the user */
	go WelcomeUser(nick, channel)

	/* Hail users */
	if _, ok := hails[nick]; ok {
		go Privmsg(channel, "All hail "+nick+"!")
	}

	log.Printf("[JOIN] %v in %v", nick, channel)

	return nil
}
