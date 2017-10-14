package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"
	"sync"
)

/*
 * irc.go
 * Handle IRC comms
 * By MagisterQuis
 * Created 20170625
 * Last Modified 20170625
 */

// Connect connects to the server
func Connect(addr string) (net.Conn, error) {
	/* Make a network connection */
	c, err := net.Dial("tcp", addr)
	if nil != err {
		return nil, err
	}
	return c, nil
}

/* doIRC logs into the IRC server, authenticates, and handles messages. Reg's
done method is called after registration */
func doIRC(c net.Conn, nick string, token []byte, reg *sync.WaitGroup) error {

	defer c.Close()

	/* Try to register and join the channel */
	if err := register(c, nick, token); nil != err {
		return err
	}
	log.Printf("Sent auth info for %v", nick)
	reg.Done()

	/* Handle received messages */
	r := textproto.NewReader(bufio.NewReader(c))
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
			if _, err := fmt.Fprintf(
				c,
				"PONG"+l[4:]+"\r\n",
			); nil != err {
				return err
			}
			continue
		}
		/* Handle every other line */
		if err := handleLine(c, l); nil != err {
			return err
		}
	}
}

/* handleLine handles individual lines */
func handleLine(w io.Writer, line string) error {
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
	switch op {
	case "001", "002", "003", "004", "353", "366", "372", "375", "376", "CAP", "MODE":
		return nil
	case "PART":
		return handlePart(w, args, src)
	case "JOIN":
		return handleJoin(w, args, src)
	case "PRIVMSG":
		return handlePrivmsg(w, src, args)
	default:
		log.Printf("Unhandled line from server: %v", line)
		return nil
	}
}

/* register registers, and joins the channel */
func register(
	c io.Writer,
	nick string,
	token []byte,
) error {
	/* Auth */
	if _, err := fmt.Fprintf(
		c,
		"PASS %s\r\n"+
			"NICK %v\r\n"+
			"USER %v quisbot quisbot :MagisterQuis' bot\r\n"+
			"CAP REQ :twitch.tv/membership\r\n",
		token,
		nick,
		nick,
	); nil != err {
		return err
	}

	return nil
}
