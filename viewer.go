package main

/*
 * viewer.go
 * DB routines handling viewers
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160831
 */

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

/* LogFirstLast logs the first and latest time a user did the given action in
the given channel.  what is an op-specific field which can be anything. */
func LogFirstLast(nick, op, what string) {
	/* Get the time as a gob */
	now := []byte(time.Now().Format(time.RFC3339))
	/* Stick it in the database */
	DB.Update(func(tx *bolt.Tx) error {
		/* Log the what and the last */
		PutTx(
			tx,
			sb("what"), sb(what),
			sb("viewers"), sb(nick), sb(op),
		)
		PutTx(
			tx,
			sb("last"), now,
			sb("viewers"), sb(nick), sb(op),
		)
		/* See if we have the first one */
		if nil != GetTx(
			tx,
			sb("first"),
			sb("viewers"), sb(nick), sb(op),
		) {
			/* We're done if we do */
			return nil
		}
		/* If not, log it */
		PutTx(
			tx,
			sb("first"), now,
			sb("viewers"), sb(nick), sb(op),
		)
		return nil
	})
	return
}

/* SetChanOp sets the chanop state of the nick in the channel */
func SetChanOp(nick, channel string, isOp bool) error {
	_, err := PutBool(
		sb(channel),
		isOp,
		sb("viewers"),
		sb(nick),
		sb("chanop"),
	)
	return err
}

/* IsChanOp returns whether the user is known to be a channel operator */
func IsChanOp(nick, channel string) bool {
	var is bool
	b := GetBool(sb(channel), sb("viewers"), sb(nick), sb("chanop"))
	if nil != b {
		is = *b
	}
	return is
}

/* ChangeAccountBalanceTx changes the viewer's bank account balance in the
transaction Tx and returns the new balance. */
func ChangeAccountBalanceTx(tx *bolt.Tx, nick string, amount int64) int64 {
	/* Change the account balance */
	bal, err := ChangeAccountBalanceBucket(
		GetBucket(tx, sb("viewers"), sb(nick)),
		amount,
	)
	if nil != err {
		panic(err.Error())
	}
	return bal
}

/* TODO: Make sure nick, globally, is small enough */

/* ChangeAccountBalanceBucket adds the value to the viewer's bank account.  The
value may be negative, to decrease the amount the user has in the bank.  The
bucket should be the viewer's bucket.  The viewer's current account balances is
returned. */
func ChangeAccountBalanceBucket(
	b *bolt.Bucket,
	amount int64,
) (cur int64, err error) {
	/* Get the current balance */
	balance := decodeBalance(b.Get(sb("credits")))
	/* Add to it */
	newbalance := balance + amount
	/* Store it */
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, newbalance)
	buf = buf[:n]
	return newbalance, b.Put(sb("credits"), buf)
}

/* GetAccountBalance gets the account balance for the viewer named nick */
func GetAccountBalance(nick string) int64 {
	var bal int64
	DB.Update(func(tx *bolt.Tx) error {
		bal = GetAccountBalanceTx(tx, nick)
		return nil
	})
	return bal
}

/* GetAccountBalanceTx gets the account balacne for the viewer in the
transaction tx. */
func GetAccountBalanceTx(tx *bolt.Tx, nick string) int64 {
	return decodeBalance(
		GetTx(tx, sb("credits"), sb("viewers"), sb(nick)),
	)
}

/* decodeBalance decodes the balance in b, or panics on error.  b may be nil,
in which case 0 is returned. */
func decodeBalance(b []byte) int64 {
	if nil == b || 0 == len(b) {
		return 0
	}
	balance, n := binary.Varint(b)
	if 0 >= n {
		log.Panicf(
			"Unable to decode account balance %q: %v",
			b,
			n,
		)
	}
	return balance
}

/* nextEventAllowedTx returns the next time this viewer's allowed to make a
new event.  If the viewer hasn't bet before, a zero time.Time is returned. */
func NextEventAllowedTx(tx *bolt.Tx, nick string) time.Time {
	/* Get the time as an RFC3339 string */
	timestamp := GetTx(tx, sb("nextbet"), sb("viewers"), sb(nick))
	if nil == timestamp {
		/* Viewer hasn't bet before */
		return time.Time{}
	}
	/* Parse the time */
	t, err := time.Parse(time.RFC3339, string(timestamp))
	if nil != err {
		panic(err.Error())
	}
	return t
}

/* SetNextEventAllowedTx sets the time the viewer may make another event. */
func SetNextEventAllowedTx(tx *bolt.Tx, nick string, next time.Time) {
	PutTx(tx, sb("nextbet"), sb(next.Format(time.RFC3339)), sb("viewers"), sb(nick))
}

/* WarnIfNotChanOp sends the user a warning if ch is a channel and the user's
not an op in that channel.  It returns true if a warning was sent */
func WarnIfNotChanOp(nick, ch string) bool {
	if IsChannel(ch) && !IsChanOp(nick, ch) {
		go Privmsg(
			ch,
			fmt.Sprintf("%v: Hey, you're not an op.", nick),
		)
		return true
	}
	return false
}
