package main

/*
 * register.go
 * Register with the twitch server
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160821
 */

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

/* Register registers, and joins the channel */
func Register(c io.Writer, username, tokenFile, channel string) error {
	/* Slurp token file */
	token, err := ioutil.ReadFile(tokenFile)
	if nil != err {
		return err
	}
	log.Printf("Read token from %v", tokenFile)

	/* Auth */
	a := fmt.Sprintf(
		"PASS %s\r\n"+
			"NICK %v\r\n"+
			"USER %v quisbot quisbot :MagisterQuis' bot\r\n"+
			"CAP REQ :twitch.tv/membership\r\n",
		token,
		username,
		username,
	)
	/* Authenticate */
	if err := Send(a); nil != err {
		return err
	}
	log.Printf("Sent authentication as %v", username)

	/* Make sure the channel has a leading # */
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}
	/* Join */
	if err := Send("JOIN " + channel); nil != err {
		return err
	}
	log.Printf("Requested to join %v", channel)

	return nil
}
