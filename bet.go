package main

/*
 * bet.go
 * Place a bet
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160821
 */

import "time"

var lastHelp time.Time

/* PlaceBet allows a user to place a bet */
func PlaceBet(nick, replyto, args string) error {
	/* If it's just help, send that back */
	if "help" == args {
		for _, s := range []string{
			"Ways to bet:",
			"!bet Ncr FOR <something> IN <time>",
			"- Bets N credts that something will happen in the " +
				"time given",
			"- Example: \"!bet 5cr FOR MagisterQuis getting 10 " +
				"kills IN 5m\"",
			"!bet Ncr AGAINST <something> IN <time>",
			"- Like FOR, but bet it won't happen",
			"!bet Ncr FOR <number>",
			"- Bet N credits that the given event will happen",
			"- Example: \"!bet 5cr FOR 301\"",
			"!bet Ncr AGAINST <number>",
			"- Bet N credits that the given event won't happen",
			"!bet account",
			"- Shows bank account balance and other such things",
			"!bet help",
			"- This help",
		} {
			if err := Privmsg(replyto, nick+": "+s); nil != err {
				return err
			}
		}
	}
	return nil
}
