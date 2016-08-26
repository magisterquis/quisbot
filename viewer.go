package main

/*
 * viewer.go
 * DB routines handling viewers
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160821
 */

import (
	"encoding/binary"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

/* BucketAndGreet returns the bucket for the viewer, and greets the viewer if
he's not been seen before.  BucketAndGreet panics on database errors. */
func BucketAndGreet(nick, channel string, tx *bolt.Tx) (*bolt.Bucket, error) {
	/* Get viewers bucket */
	vb, err := tx.CreateBucketIfNotExists(sb("viewers"))
	if nil != err {
		log.Panicf(
			"Unable to get/create viewers bucket: %v",
			err,
		)
	}
	/* If we have a bucket, it's not the first time */
	if b := vb.Bucket(sb(nick)); nil != b {
		return b, nil
	}
	/* It's the first time, so send a welcome message */
	if err := WelcomeUser(nick, channel); nil != err {
		return nil, err
	}
	/* Make the user a bucket */
	b, err := vb.CreateBucket(sb(nick))
	if nil != err {
		log.Panicf(
			"Unable to create bucket for %q: %v",
			sb(nick),
			err,
		)
	}
	return b, nil
}

/* LogFirstLast logs the first and latest time a user did the given action in
the given channel.  what is an op-specific field which can be anything. */
func LogFirstLast(nick, op, replyto, what string) {
	/* Get the time as a gob */
	now := []byte(time.Now().Format(time.RFC3339))
	/* Stick it in the database */
	if derr := DB.Update(func(tx *bolt.Tx) error {
		/* Get the user's bucket, maybe welcome him? */
		ub, err := BucketAndGreet(nick, replyto, tx)
		if nil != err {
			return err
		}
		/* Get the user's op bucket */
		ob, err := ub.CreateBucketIfNotExists(sb(op))
		if nil != err {
			return err
		}
		/* If there's no first, this is the first. */
		if nil == ob.Get(sb("first")) {
			if err := ob.Put(sb("first"), now); nil != err {
				return err
			}
		}
		if err := ob.Put(sb("last"), now); nil != err {
			return err
		}
		if err := ob.Put(sb("what"), sb(what)); nil != err {
			return err
		}
		return nil
	}); nil != derr {
		log.Printf(
			"Unable to log first/last %v for %v (from %v): %v",
			nick,
			op,
			replyto,
			derr,
		)
	}

	return
}

/* SetChanOp sets the chanop state of the nick in the channel */
func SetChanOp(nick, channel string, isOp bool) error {
	/* Stick it in the database */
	if derr := DB.Update(func(tx *bolt.Tx) error {
		/* Get the viewer's bucket, maybe welcome him? */
		vb, err := BucketAndGreet(nick, channel, tx)
		if nil != err {
			return err
		}
		/* Get the viewer's chanop bucket */
		ob, err := vb.CreateBucketIfNotExists(sb("chanop"))
		if nil != err {
			return err
		}
		/* Set the user's chanop status */
		return PutBool(ob, sb(channel), isOp)
	}); nil != derr {
		return derr
	}
	return nil

}

/* IsChanOp returns whether the user is known to be a channel operator */
func IsChanOp(nick, channel string) (bool, error) {
	var isOp bool /* Opness to return */
	if derr := DB.Update(func(tx *bolt.Tx) error {
		/* Viewer's bucket */
		vb, err := BucketAndGreet(nick, channel, tx)
		if nil != err {
			return err
		}
		/* Viewer's chanop bucket */
		ob, err := vb.CreateBucketIfNotExists(sb("chanop"))
		if nil != err {
			return err
		}
		isOp, _, err = GetBool(ob, sb(channel))
		if nil != err {
			return err
		}
		return nil
	}); nil != derr {
		return false, derr
	}
	return isOp, nil
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
	b, err := Get([][]byte{sb("viewers"), sb(nick)}, sb("credits"))
	return decodeBalance(b), err
}

/* decodeBalance decodes the balance in b, or panics on error.  b may be nil,
in which case 0 is returned. */
func decodeBalance(b []byte) int64 {
	if nil == b {
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
