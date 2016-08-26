package main

/*
 * send.go
 * Send various messages
 * By MagisterQuis
 * Created 20160820
 * Last Modified 20160820
 */

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

/* SENDDELAY is the time to sleep between sends */
const SENDDELAY = 1500 * time.Millisecond

/* sendLock keeps sends thread-safe */
var sendLock = &sync.Mutex{}

/* send sends the message, as long as there's a connection */
func Send(msg string) error {
	/* Make sure there's a connection */
	if nil == ServerConnection {
		return fmt.Errorf("No connection yet to send %q", msg)
	}

	sendLock.Lock()
	defer sendLock.Unlock()

	/* CRLF-terminate message */
	if !strings.HasSuffix(msg, "\r\n") {
		msg += "\r\n"
	}

	/* Make sure message isn't too long */
	b := []byte(msg)
	if 512 < len(msg) {
		return fmt.Errorf("message %q too long (%v > 512)", msg, len(b))
	}

	/* Try to send it */
	_, err := ServerConnection.Write(b)
	if nil != err {
		return err
	}

	time.Sleep(SENDDELAY)
	return nil
}

/* Privmsg sends a private message to the target.  If the message is too long,
it will be split up.  Message may be a printf-like format string, filled in
by args. */
func Privmsg(target, message string, args ...interface{}) error {
	/* Fill out the message */
	message = fmt.Sprintf(message, args...)
	/* Per-chunk leader */
	cmd := fmt.Sprintf("PRIVMSG %v :", target)
	if 510 <= len(cmd) {
		return fmt.Errorf("Nick too long for command (>510): %v", cmd)
	}

	/* Stream from which to read a long message */
	stream := bytes.NewBufferString(message)

	/* Buffer for message chunk */
	buf := make([]byte, 510-len(cmd))
	go func() {
		/* Read and send until we have no buf left */
		for 0 != stream.Len() {
			buf = buf[:cap(buf)]
			/* Read a chunk */
			n, err := stream.Read(buf)
			if 0 == n {
				if nil != err && io.EOF != err {
					log.Printf(
						"Unable to read from buffer "+
							"sending %v to %v: %v",
						message,
						target,
						err,
					)
				}
				return
			}
			buf = buf[:n]
			/* Send it */
			if serr := Send(
				cmd + string(buf),
			); nil != serr {
				log.Printf("Error sending %b: %v", buf, err)
				return
			}
		}
	}()
	return nil
}
