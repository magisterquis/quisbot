package main

/*
 * receive.go
 * Handle received messages
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160821
 */

import (
	"bufio"
	"fmt"
	"log"
	"net/textproto"
	"strings"
)

/* HandleRX handles received lines */
func HandleRX() error {
	r := textproto.NewReader(bufio.NewReader(ServerConnection))
	for {
		/* Grab a line */
		l, err := r.ReadLine()
		if nil != err {
			return err
		}
		/* Ignore blank lines */
		if "" == l {
			continue
		}
		/* Respond if it's a ping */
		if strings.HasPrefix(strings.ToLower(l), "ping") {
			if err := Send("PONG" + l[4:]); nil != err {
				return err
			}
			continue
		}
		/* Handle every other line */
		if err := handleLine(l); nil != err {
			return err
		}
	}
	return nil
}

/* handleLine handles individual lines */
func handleLine(line string) error {
	/* Split message into parts */
	parts := strings.SplitN(line, " ", 3)
	if 0 == len(parts) {
		return fmt.Errorf("No parts to message %q", line)
	}
	var src, op, args string

	/* Name the parts */
	switch len(parts) {
	case 3:
		args = parts[2]
		fallthrough
	case 2:
		op = parts[1]
		fallthrough
	case 1:
		src = strings.SplitN(parts[0][1:], "!", 2)[0]
	}

	/* Function to call for this line */
	var f func(src, op, args string) error
	switch op {
	case "CAP":
		f = logLine
	case "001", "002", "003", "004", "353", "366", "372", "375", "376":
		f = ignoreLine
	case "JOIN":
		f = HandleJoin
	case "PRIVMSG":
		f = HandlePrivmsg
	case "MODE":
		f = HandleMode
		/* TODO: Handle quits */
	}

	/* Print lines we can't handle */
	if nil == f {
		log.Printf("Unhandled line: %v", line)
		return nil
	}

	/* Call the message-specific function */
	return f(src, op, args)
}

/* ignoreLine is a no-op to ignore a line that's expected */
func ignoreLine(string, string, string) error { return nil }

/* logLine logs the line */
func logLine(src, op, args string) error {
	log.Printf("[PROTOCOL] %v from %v: %v", op, src, args)
	return nil
}

/* TODO: Handle parts */
//2016/08/24 05:33:47.026705 Unhandled line: :moobot!moobot@moobot.tmi.twitch.tv PART #magisterquis
