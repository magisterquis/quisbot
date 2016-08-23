package main

/*
 * payroll.go
 * Give active people credits
 * By MagisterQuis
 * Created 20160823
 * Last Modified 20160823
 */

import (
	"bytes"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

const (
	/* PAYRATE is the amount of credits to hand out every interval */
	PAYRATE = 10
	/* PAYINTERVAL is how often viewers get paid */
	PAYINTERVAL = 15 * time.Minute
)

/* Payroll gives active people credits every hour */
func Payroll() {
	/* Every second, check if it's a new interval */
	last := time.Now().Truncate(PAYINTERVAL)
	for {
		now := time.Now().Truncate(PAYINTERVAL)
		if !now.Equal(last) {
			/* When we've crossed into a new interval, pay out */
			go payOut(
				[]byte(last.Format(time.RFC3339)),
				[]byte(now.Format(time.RFC3339)),
			)
		}
		last = now
		time.Sleep(time.Second)
	}
}

/* payOut pays people active between the two times [from, to), which should be
a byte slice containing an RFC3339 time.  */
func payOut(from, to []byte) {
	nPaid := 0
	if derr := DB.Update(func(tx *bolt.Tx) error {
		/* Viewers bucket */
		vb, err := tx.CreateBucketIfNotExist(sb("viewers"))
		if nil != err {
			return err
		}
		/* For each viewer, increase their pay if they've talked
		recently */
		c := vb.Cursor()
		for k, v := c.First(); nil != k; k, v = c.Next() {
			/* If the key's value is non-nil, the database has got
			corrupt. */
			if nil != v {
				log.Fatalf(
					"Database corruption detected.  "+
						"Viewer key %q has value %q.",
					k,
					v,
				)
			}
			/* Pay the viewer */
			b := vb.Bucket(k)
			if nil == b {
				log.Fatalf("Unable to get viewer bucket %q", k)
			}
			if err := Pay(string(k), b, from, to); nil != err {
				return err
			}
		}
	}); nil != derr {
		log.Printf("Unable to hand out pay: %v", err)
	}
}

/* Pay pays the viewer v with bucket b if he's been seen in the given
time interval [from, to). */
func Pay(v string, b *bolt.Bucket, from, to []byte) error {
	/* Last PRIVMSG time */
	pb := b.Bucket(sb("PRIVMSG"))
	if nil == pb {
		log.Printf("%v has no privmsg bucket", v) /* DEBUG */
		return nil
	}
	last := pb.Get(sb("last"))
	if nil == last {
		log.Printf("%v has no privmsg last", v) /* DEBUG */
		return nil
	}
	log.Printf("%v privmsg last: %s", v, last)
	/* Make sure the user was active in the right time interval */
	if (-1 == bytes.Compare(last, from)) ||
		!(-1 == bytes.Compare(last, to)) {
		log.Printf("Not paying %v", v)
		return nil
	}
	/* Increase the bank account by the proper number of credits */
	if err := AccountBucketChange(b, PAYRATE); nil != err {
		return err
	}
	log.Printf("[PAYROLL] Paid %v to %v (%s)", PAYRATE, v, last)
	return nil
}
