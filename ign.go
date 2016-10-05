package main

/*
 * ign.go
 * Remember users' IGNs
 * By MagisterQuis
 * Created 20161004
 * Last Modified 20161004
 */

import (
	"log"
	"strings"
)

const IGN = "magisterquis"

/* IGN either prints the default IGN or the request user's IGN */
func IGNCommand(nick, replyto, args string) error {
	parts := strings.Fields(args)

	/* Request for the baked-in IGN */
	if 0 == len(parts) {
		log.Printf("[IGN] Request for !ign by %v", nick)
		go Privmsg(replyto, "%v: His IGN is %v", nick, IGN)
		return nil
	}

	/* Get the IGN for a user */
	if 1 == len(parts) {
		ign := GetIGN(parts[0])
		/* Don't have one */
		if "" == ign {
			log.Printf(
				"[IGN] Don't have an IGN for %v "+
					"(requested by %v)",
				parts[0],
				nick,
			)
			go Privmsg(
				replyto,
				"%v: Sorry, I don't have an IGN for %v",
				nick,
				parts[0],
			)
			return nil
		}
		log.Printf("[IGN] %v -> %v", parts[0], ign)
		go Privmsg(
			replyto,
			"%v: The IGN I have for %v is %v",
			nick,
			parts[0],
			ign,
		)
		return nil
	}

	/* Set the IGN for a user, if it's a chanop asking */
	if WarnIfNotChanOp(nick, replyto) {
		return nil
	}

	/* Try to set the IGN */
	SetIGN(parts[0], parts[1])
	log.Printf("[IGN] Set IGN for %v to %v", parts[0], parts[1])
	go Privmsg(
		replyto,
		"%v: Saved IGN %v for %v",
		nick,
		parts[1],
		parts[0],
	)
	return nil
}
