package main

/*
 * bet.go
 * Place a bet
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160827
 */

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	BETHELPWAIT       = 2 * time.Minute  /* Time between !bet helps */
	MINEVENTLEN       = 2 * time.Minute  /* Minimum bet event length */
	MAXEVENTLEN       = 15 * time.Minute /* Maximum bet event length */
	MAXCONCURRENTBETS = byte(5)          /* Maximum number of live bets */
)

var (
	lastHelp time.Time             /* Last time help message was printed */
	eventre  = regexp.MustCompile( /* New event */
		`^\s*((?:0x)?\d+)\s+(?i:THAT)\s+(.*)\s+(?i:IN)\s+(\S+)\s*$`,
	)
	betre = regexp.MustCompile( /* Betting on an event */
		`^\s*((?:0x)?\d+)\s+((?i)FOR|AGAINST)\s+(\d+)\s*$`,
	)
	bethelprl = NewRateLimiter(BETHELPWAIT)
)

/* Timers related to events */
var betTimers = [256]struct {
	bet   *time.Timer /* Betting's over */
	event *time.Timer /* The event's done */
}{}

/* PlaceBet allows a user to place a bet */
func PlaceBet(nick, replyto, args string) error {
	answered := true /* Will be true if it's an easy match */
	/* Single-word commands */
	switch args {
	case "help":
		/* If it's just help, send that back */
		go sendBetHelp(nick, replyto)
	case "ophelp":
		/* Help for ops */
		go sendBetOpHelp(nick, replyto)
	case "list":
		go listEvents(nick, replyto)
	/* TODO: Finish this */
	default:
		answered = false
	}
	if answered {
		return nil
	}

	/* Check if someone's putting money on something */
	if parts := eventre.FindStringSubmatch(args); nil != parts {
		go addEvent(nick, replyto, parts[1], parts[2], parts[3])
	} else if parts := betre.FindStringSubmatch(args); nil != parts {
		/* Maybe someone's placing a bet */
		go makeBet(nick, replyto, parts[1], parts[2], parts[3])
	} else if strings.HasPrefix(args, "resettimer") {
		/* Reset someone's bet timer */
		go resetNextEventTimer(nick, replyto, args[10:])
	} else if strings.HasPrefix(args, "killbet") {
		/* Remove a bet from the list */
		go killBet(nick, replyto, args[7:])
	} else {

		go Privmsg(replyto, fmt.Sprintf(
			"%v: I don't understand that.  Please try !bet help.",
			nick,
		))
	}

	/* TODO: Make error messages much friendier. */
	return nil
}

/* sendBetOpHelp lists the commands available to ops */
func sendBetOpHelp(nick, replyto string) {
	if WarnIfNotChanOp(nick, replyto) {
		return
	}
	go Privmsg(replyto, fmt.Sprintf("%v: !bet {resettimer,killbet}", nick))
}

/* sendBetHelp sends the long help output to the requestor */
func sendBetHelp(nick, replyto string) {
	/* Don't send too often */
	if 0 != bethelprl.Until() {
		/* Tell the user */
		go Privmsg(replyto, fmt.Sprintf(
			"%v: Sorry, I can't help that fast.",
			nick,
		))
		return
	}
	/* Send command help */
	for _, s := range []string{
		"Betting help:",
		"!bet N THAT <something> IN <time>",
		"- Bets N internets that something will happen in the " +
			"time given",
		"- Example: \"!bet 5 THAT MagisterQuis gets 10 kills IN 5m\"",
		"!bet N FOR <ID>",
		"!bet N AGAINST <ID>",
		"- Bet on whether event ID will happen or not",
		"!bet list",
		"- List the events up for betting",
		"!bet help",
		"- This help",
	} {
		Privmsg(replyto, nick+": "+s)
	}
	lastHelp = time.Now()
}

/* TODO: a way to cancel events */
/* TODO: a way to list events */

/* parseBetAmount parses the bet as a string, ensures it's over 0, and returns
the value. */
func parseBetAmount(bet string) (int64, error) {
	/* Convert the bet to a number */
	bal, err := strconv.ParseInt(bet, 0, 64)
	if nil != err {
		return bal, err
	}
	if 0 >= bal {
		return bal, fmt.Errorf("a bet cannot be negative")
	}
	return bal, nil
}

/* addEvent adds the event to be bet upon, expiring in time.  An initial bet by
nick of value bet will be placed. */
func addEvent(nick, replyto, bet, event, duration string) {
	/* TODO: Chop function up a bit */
	var (
		betID    byte /* Bet ID */
		eventEnd time.Time
		betEnd   time.Time
		now      time.Time
	)
	if derr := DB.Update(func(tx *bolt.Tx) error {
		now = time.Now()
		/* Work out the bet */
		betVal, err := parseBetAmount(bet)
		if nil != err {
			return err
		}

		/* Make sure the viewer has the funds */
		if GetAccountBalanceTx(tx, nick) < betVal {
			return fmt.Errorf("insufficent funds")
		}

		/* Work out the time the event expires, and the last time bets
		are accepted */
		betDuration, err := time.ParseDuration(duration)
		if nil != err {
			return err
		}
		if MAXEVENTLEN < betDuration {
			return fmt.Errorf(
				"it'd take too much time (>%v)",
				MAXEVENTLEN,
			)
		}
		if MINEVENTLEN > betDuration {
			return fmt.Errorf(
				"it's not enough time (<%v)",
				MINEVENTLEN,
			)
		}
		eventEnd = now.Add(betDuration)
		betEnd = now.Add(betDuration / 2)

		/* Make sure the viewer can make an event */
		nextTime := NextEventAllowedTx(tx, nick)
		if time.Now().Before(nextTime) {
			return fmt.Errorf(
				"not allowed until %v",
				nextTime,
			)
		}

		/* TODO: Have a way to reset this */

		/* Get the bet bucket */
		betBucket := GetBucket(tx, sb("bets"))
		if nil != err {
			panic(err.Error())
		}

		/* TODO: Log the bet has been placed */

		/* Get the current state of the bets */
		c := betBucket.Cursor()
		n := byte(0) /* Number of bets in play */
		for k, v := c.First(); nil != k; k, v = c.Next() {
			n++
			if 0 == len(k) {
				panic("0-length bet bucket key")
			}
			/* Don't use this ID if it's in use */
			if betID == k[0] {
				if 0xFF == betID {
					return fmt.Errorf(
						"too many bets (>%v)",
						betID,
					)
				}
				betID++
			}
			/* If it's not a bucket, it's an error */
			if nil != v {
				log.Panicf("%q not a bet bucket", k)
			}
			b := betBucket.Bucket(k)
			if nil == b {
				log.Panicf("bet bucket %q doesn't exist", k)
			}
			p := b.Get(sb("event"))
			if nil == p {
				log.Panicf("bet %q has empty proposition", k)
			}
			if 0 == bytes.Compare(p, sb(event)) {
				return fmt.Errorf("already submitted")
			}
		}

		/* Make sure there's not too many already */
		if n >= MAXCONCURRENTBETS {
			return fmt.Errorf(
				"only %v bets allowed at once",
				MAXCONCURRENTBETS,
			)
		}

		/* Add the bet to the list */
		b, err := betBucket.CreateBucketIfNotExists([]byte{betID})
		if nil != err {
			return err
		}
		for k, v := range map[string]string{
			"event":   event,
			"end":     eventEnd.Format(time.RFC3339),
			"lastbet": betEnd.Format(time.RFC3339),
		} {
			if err := b.Put(sb(k), sb(v)); nil != err {
				return err
			}
		}

		/* The bettor makes the first wager */
		if err := addBetTx(tx, nick, betVal, betID, true); nil != err {
			return err
		}

		/* Note that the viewer can't bet until this one's done */
		SetNextEventAllowedTx(tx, nick, eventEnd)

		return nil
	}); nil != derr {
		go Privmsg( /* TODO: Make much nicer */
			replyto,
			"%v: Sorry, you can't make that bet: %v",
			nick,
			derr,
		)
		return
	}
	/* Fire off goroutines to let channel know betting's over and
	to ask whether the event happened */
	go Privmsg( /* TODO: Make nicer */
		replyto,
		"Bet's in!  You have %v to place your bets for event %v: %v",
		roundToSeconds(betEnd.Sub(time.Now())),
		betID,
		event,
	)
	/* Start timers to notify the channel that the betting and event are
	over */
	/* TODO: Implement the below */
	betTimers[betID].bet = time.AfterFunc(
		betEnd.Sub(now),
		func() {
			betFinished(replyto, betID)
		},
	)
	betTimers[betID].event = time.AfterFunc(
		eventEnd.Sub(now),
		func() {
			eventFinished(replyto, betID)
		},
	)
	log.Printf(
		"[BET] %v (%v) by %v in %v BetEnd:%v EventEnd:%v",
		event,
		betID,
		nick,
		replyto,
		betEnd,
		eventEnd,
	)

	//eventEnd := now.Add(betDuration)
	//betEnd := now.Add(betDuration / 2)
	/* TODO: Tell user it all worked */
}

/* addBetTx adds a bet of n to the event with the given id for the viewer with
the given nick in the given transaction.  Whether the bet is added for or
against is controlled by isFor.  The viewer's funds are checked to make sure
they are sufficient. */
func addBetTx(
	tx *bolt.Tx,
	nick string,
	n int64,
	betID byte,
	isFor bool,
) error {
	/* Make sure bet isn't negative */
	if 0 > n {
		return fmt.Errorf("negative bets not allowed")
	}
	/* Make sure viewer has enough */
	if n > GetAccountBalanceTx(tx, nick) {
		return fmt.Errorf("you don't have enough money")
	}
	/* Get bucket for bet, for or against */
	var way string
	if isFor {
		way = "for"
	} else {
		way = "against"
	}
	b := GetBucket(tx, sb("bets"), []byte{betID}, sb(way))
	/* Get previous bet for bettor */
	var prev int64
	var nr int
	prevBuf := b.Get(sb(nick))
	if nil != prevBuf {
		prev, nr = binary.Varint(prevBuf)
		if 0 > nr {
			log.Panicf("Unvarint %q failed: %v", prevBuf, nr)
		}
	}
	/* Subtract the bet (update) from the bettor's bank account */
	ChangeAccountBalanceTx(tx, nick, -1*n)
	/* Add new amount to it */
	bet := prev + n
	/* Stick it back in */
	vbuf := make([]byte, binary.MaxVarintLen64)
	vlen := binary.PutVarint(vbuf, bet)
	vbuf = vbuf[:vlen]
	if err := b.Put(sb(nick), vbuf); nil != err {
		panic(err.Error())
	}
	return nil
}

/* makeBet puts bet monies on (or against, via way) event id */
func makeBet(nick, replyto, bet, way, id string) {
	if derr := DB.Update(func(tx *bolt.Tx) error {
		/* Work out bet amount */
		amt, err := parseBetAmount(bet)
		if nil != err {
			return err
		}
		/* Lower-case the way */
		way = strings.ToLower(way)
		/* Bet ID as a byte */
		bint, err := strconv.ParseUint(bet, 0, 8)
		if nil != err {
			return err
		}
		bid := byte(bint)
		/* TODO: Use ParseUint better when there should be no negative ints */
		/* Bet bucket */
		b := GetBucket(tx, sb("bets"))
		/* Make sure bet ID exists */
		b = b.Bucket([]byte{bid})
		if nil == b {
			return fmt.Errorf("event %v does not exist", bid)
		}
		/* Make sure betting time hasn't expired */
		lb := b.Get(sb("lastbet"))
		if nil == lb {
			log.Panicf("Event %v has no lastbet", bid)
		}
		if -1 != bytes.Compare(
			sb(time.Now().Format(time.RFC3339)),
			lb,
		) {
			return fmt.Errorf("betting has closed")
		}

		/* Add bet to event */
		var isFor bool
		switch way {
		case "for":
			isFor = true
		case "against":
			isFor = false
		default:
			return fmt.Errorf(
				"a bet must be FOR or AGAINST, not %v",
				way,
			)
		}
		return addBetTx(tx, nick, amt, bid, isFor)
	}); nil != derr {
		go Privmsg(replyto, fmt.Sprintf(
			"%v: error betting %v %v: %v",
			nick,
			way,
			id,
			derr,
		))
		return
	}
	go Privmsg(replyto, fmt.Sprintf(
		"%v: bet %v %v %v",
		nick,
		bet,
		way,
		id,
	))
}

/* betFinished is called when betting is closed for the bet betID.  It will let
the channel (or whatever replyto is) know. */
func betFinished(replyto string, betID byte) {
	/* TODO: Implement this */
	panic("implement this")
}

/* eventFinished is called when betting is closed for the bet betID.  It will
let the channel (or whatever replyto is) know. */
func eventFinished(replyto string, betID byte) {
	/* TODO: Implement this */
	panic("implement this")
}

/* resetNextEventTimer reset's the target nick's bet-placing timer */
func resetNextEventTimer(nick, replyto, targetList string) {
	/* Split nicks into a list */
	ns := strings.Split(targetList, " ")
	DB.Update(func(tx *bolt.Tx) error {
		/* Reset each nick's timer */
		for _, n := range ns {
			n = strings.TrimSpace(n)
			if "" == n {
				continue
			}
			/* TODO: Fix timing code */
			SetNextEventAllowedTx(tx, n, time.Time{})
			log.Printf(
				"[BET] Reset next event allowed time for %v",
				n,
			)
			go Privmsg(replyto, fmt.Sprintf(
				"%v: Reset next event timer for %v",
				nick,
				n,
			))
		}
		return nil
	})
	return
}

/* listEvents lists the bettable events */
func listEvents(nick, replyto string) {
	eventListings := make([]string, 0)

	DB.Update(func(tx *bolt.Tx) error {
		now := time.Now()
		b := GetBucket(tx, sb("bets"))
		/* Iterate over the current events, stick them in es */
		c := b.Cursor()
		for k, _ := c.First(); nil != k; k, _ = c.Next() {
			if 1 != len(k) {
				log.Panicf("Invalid bet key: %q", k)
			}
			/* Bucket specific to this event */
			eb := b.Bucket(k)
			eventListings = append(
				eventListings,
				getEventListingFromBucket(k[0], eb, now),
			)
		}
		return nil
	})
	go func() {
		Privmsg(replyto, nick+": Current events...")
		for _, e := range eventListings {
			Privmsg(replyto, e)
		}
		log.Printf("[BET] Listed events for %v in %v", nick, replyto)
	}()
}

/* getEventListingFromBucket gets a one-liner event listing from the event
stored in the bucket b */
func getEventListingFromBucket(bid byte, b *bolt.Bucket, now time.Time) string {
	var (
		lastbet      time.Duration
		eventend     time.Duration
		eventListing string
	)
	/* Get time betting ends */
	t, err := time.Parse(time.RFC3339, bs(b.Get(sb("lastbet"))))
	if nil != err {
		panic(err.Error())
	}
	if t.Before(now) {
		lastbet = t.Sub(now)
	}
	/* Get the time until event ends */
	t, err = time.Parse(time.RFC3339, bs(b.Get(sb("end"))))
	if nil != err {
		panic(err.Error())
	}
	if t.Before(now) {
		eventend = t.Sub(now)
	}
	/* Get event */
	event := bs(b.Get(sb("event")))
	/* Get amounts bet */
	bb, err := b.CreateBucketIfNotExists(sb("for"))
	if nil != err {
		panic(err.Error())
	}
	forBet := totBets(bb)
	bb, err = b.CreateBucketIfNotExists(sb("against"))
	if nil != err {
		panic(err.Error())
	}
	againstBet := totBets(bb)

	/* Start with the bet ID */
	eventListing = fmt.Sprintf("%d: ", bid)

	/* Note the time remaining until the event */
	if 0 >= eventend {
		eventListing += "[DONE]"
	} else {
		tr := roundToSeconds(eventend)
		if 0 == tr {
			tr = eventend
		}
		eventListing += fmt.Sprintf("[%v]", tr)
	}
	/* Add the event and the amount bet for and against */
	eventListing += fmt.Sprintf(
		" %v (%vF/%vA)",
		event,
		forBet,
		againstBet,
	)
	/* If there's stil time to bet, add that */
	if 0 < lastbet {
		tr := roundToSeconds(lastbet)
		if 0 == tr {
			tr = lastbet
		}
		eventListing += fmt.Sprintf(" - %v left to bet", tr)
	}
	return eventListing
}

/* totBets sums up the bets in the given for/against bucket */
func totBets(b *bolt.Bucket) int64 {
	var tot int64
	b.ForEach(func(k, v []byte) error {
		/* Unvarint v */
		i, n := binary.Varint(v)
		if 0 >= n {
			panic("Failed to decode varint %q totalling bets")
		}
		tot += i
		return nil
	})
	return tot
}

/* roundToSeconds rounds d to the nearest second */
func roundToSeconds(d time.Duration) time.Duration {
	return time.Duration(int64(d.Seconds())) * time.Second
}

/* killBet removes a bet from the list */
func killBet(nick, replyto, args string) {
	/* If nick's not an op, not allowed */
	if WarnIfNotChanOp(nick, replyto) {
		return
	}

	/* Get bet numbers to kill */
	ns := strings.Split(args, " ")
	DB.Update(func(tx *bolt.Tx) error {
		/* Get the bucket for the bets */
		b := GetBucket(tx, sb("bets"))
		/* Kill each bet by ID */
		for _, n := range ns {
			if "" == n {
				continue
			}
			/* Convert each number to a byte */
			bint, err := strconv.ParseUint(n, 0, 8)
			if nil != err {
				go Privmsg(
					replyto,
					"%v: Could not parse %q: %v",
					nick,
					n,
					err,
				)
				continue
			}
			bid := byte(bint)
			/* Remove it */
			if err := b.DeleteBucket([]byte{bid}); nil != err {
				/* Replace bolt message with friendlier error
				message */
				if bolt.ErrBucketNotFound == err {
					err = fmt.Errorf(
						"does not exist",
					)
				}
				go Privmsg(
					replyto,
					"%v: Unable to kill %v: %v",
					nick,
					bid,
					err,
				)
				continue
			}
			go Privmsg(replyto, "%v: Killed event %v", nick, bid)
			log.Printf(
				"[BET] Killed event %v as per %v in %v",
				bid,
				nick,
				replyto,
			)
		}
		return nil
	})
}
