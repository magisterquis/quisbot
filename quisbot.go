package main

/*
 * quisbot.go
 * Twitch IRC bot
 * By MagisterQuis
 * Created 20160820
 * Last Modified 20160821
 */

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

/* Interval before reconnecting */
const RECONINT = 30 * time.Second

var (
	ServerConnection net.Conn /* Connection to server */
	DB               *bolt.DB
)

func main() {
	var (
		user = flag.String(
			"u",
			"quisbot",
			"Twitch `username`",
		)
		tokenFile = flag.String(
			"t",
			"./token",
			"Name of `file` from which to read oauth token",
		)
		channel = flag.String(
			"c",
			"magisterquis",
			"Twitch `channel` to join",
		)
		logFile = flag.String(
			"l",
			"quisbot.log.txt",
			"Name of `logfile`",
		)
		dbFile = flag.String(
			"db",
			"./quisbot.db",
			"Database `file`",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %v [options]

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Log to a file and stderr */
	lf, err := os.OpenFile(
		*logFile,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if nil != err {
		log.Fatalf("Unable to open logfile %v: %v", *logFile, err)
	}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(io.MultiWriter(os.Stdout, lf))

	/* Open the database */
	DB, err = bolt.Open(*dbFile, 0600, nil)
	if nil != err {
		log.Fatalf("Unable to open database %v: %v", *dbFile, err)
	}
	log.Printf("Opened database in %v", DB.Path())

	/* Trap Ctrl+C, etc, to close DB */
	CatchInt()

	/* Maintain a connection to twitch */
	for {
		/* Make a network connection */
		ServerConnection, err = net.Dial(
			"tcp",
			"irc.chat.twitch.tv:6667",
		)
		if nil != err {
			log.Printf("Unable to connect to twitch: %v", err)
			goto SLEEP
		}
		log.Printf("Connected to %v", ServerConnection.RemoteAddr())
		/* Register */
		if err = Register(
			ServerConnection,
			*user,
			*tokenFile,
			*channel,
		); nil != err {
			log.Printf("Unable to register: %v", err)
			goto SLEEP
		}
		if err = HandleRX(); nil != err {
			log.Printf("RX error: %v", err)
			goto SLEEP
		}
	SLEEP:
		time.Sleep(RECONINT)

	}
}