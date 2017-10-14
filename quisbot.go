// quisbot is a small IRC bot which welcomes users
package main

/*
 * quisbot.go
 * Second try at a friendly IRC bot
 * By MagisterQuis
 * Created 20170625
 * Last Modified 20170625
 */

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	var (
		nick = flag.String(
			"nick",
			"quisbot",
			"IRC `nickname`",
		)
		configFile = flag.String(
			"conf",
			"quisbot.conf",
			"Config `file`",
		)
		tokenFile = flag.String(
			"token",
			"quisbot.token",
			"OAuth token `file`",
		)
		pollInterval = flag.Duration(
			"poll",
			time.Second,
			"Configuration change poll `interval`",
		)
		addr = flag.String(
			"addr",
			"irc.chat.twitch.tv:6667",
			"IRC server `address`",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %v [options]

Sits in channels and welcomes people.  Polls the config file for changes.

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(os.Stdout)

	/* Read in token file */
	token, err := ioutil.ReadFile(*tokenFile)
	if nil != err {
		log.Fatalf("Unable to read token from %v: %v", *tokenFile, err)
	}
	log.Printf("Read token from %v", *tokenFile)

	/* Try to connect to the server */
	c, err := Connect(*addr)
	if nil != err {
		log.Fatalf("Unable to connect to %v: %v", *addr, err)
	}
	log.Printf("Connected to %v", c.RemoteAddr())

	/* Read config */
	ready := &sync.WaitGroup{}
	ready.Add(1)
	go startConfigPoll(*configFile, *pollInterval, ready, c)

	/* Connect to server, greet users.  Random disconnects happen. */
	for {
		if err := doIRC(c, *nick, token, ready); nil != err {
			log.Printf("IRC Error: %v", err)
		}
	}
}
