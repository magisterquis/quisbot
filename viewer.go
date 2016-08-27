package main

/*
 * viewer.go
 * DB routines handling viewers
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160826
 */

import (
	"encoding/binary"
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
	if derr := DB.Update(func(tx *bolt.Tx) error {
		/* Log the what and the last */
		if err := PutTx(
			tx,
			sb("what"), sb(what),
			sb("viewers"), sb(nick), sb(op),
		); nil != err {
			return err
		}
		if err := PutTx(
			tx,
			sb("last"), now,
			sb("viewers"), sb(nick), sb(op),
		); nil != err {
			return err
		}
		/* See if we have the first one */
		f, err := GetTx(
			tx,
			sb("first"),
			sb("viewers"), sb(nick), sb(op),
		)
		if nil != err {
			return err
		}
		/* We're done if we do */
		if nil != f {
			return nil
		}
		/* If not, log it */
		if err := PutTx(
			tx,
			sb("first"), now,
			sb("viewers"), sb(nick), sb(op),
		); nil != err {
			return err
		}
		return nil
	}); nil != derr {
		log.Printf(
			"Unable to log first/last %v for %v: %v",
			nick,
			op,
			derr,
		)
	}

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
func IsChanOp(nick, channel string) (bool, error) {
	var is bool
	b, err := GetBool(sb(channel), sb("viewers"), sb(nick), sb("chanop"))
	if nil != b {
		is = *b
	}
	return is, err
}

/* ChangeAccountBucket adds the value to the viewer's bank account.  The value
may be negative, to decrease the amount the user has in the bank.  The bucket
should be the viewer's bucket.  The viewer's current account balances is
returned. */
func ChangeAccountBalance(
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
func GetAccountBalance(nick string) (int64, error) {
	b, err := Get(sb("credits"), sb("viewers"), sb(nick))
	return decodeBalance(b), err
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

/* TODO: Functions to get, set, and change credit balance */
/* TODO: Split change into functions taking a nick and another taking a
bucket */
