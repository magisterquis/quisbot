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
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	/* PAYRATE is the amount of credits to hand out every interval */
	PAYRATE = 10
	/* PAYINTERVAL is how often viewers get paid */
	PAYINTERVAL = 15 * time.Minute
	//PAYINTERVAL   = 10 * time.Second
	CURRENCYUNITS = "n"
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
	paid := map[string]string{} /* Who got paid and where to tell people */
	if derr := DB.Update(func(tx *bolt.Tx) error {
		/* Viewers bucket */
		vb, err := tx.CreateBucketIfNotExists(sb("viewers"))
		if nil != err {
			return err
		}
		/* For each viewer, increase their pay if they've talked
		recently */
		c := vb.Cursor()
		for k, v := c.First(); nil != k; k, v = c.Next() {
			/* Ignore bots */
			if _, ok := BOTS[bs(k)]; ok {
				continue
			}
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
			bal, replyto, err := Pay(string(k), b, from, to)
			if nil != err {
				return err
			}
			/* Not the payment */
			if _, ok := paid[replyto]; !ok {
				paid[replyto] = ""
			}
			paid[replyto] += fmt.Sprintf(" %v (%v).", bs(k), bal)
		}
		return nil
	}); nil != derr {
		log.Printf("Unable to hand out pay: %v", derr)
	}
	for k, v := range paid {
		go Privmsg(k, "[PAYROLL]"+v)
	}
}

/* Pay pays the viewer v with bucket b if he's been seen in the given
time interval [from, to).  The viewer's account balance, replyto, and any
errors are returned. */
func Pay(v string, b *bolt.Bucket, from, to []byte) (int64, string, error) {
	/* Last PRIVMSG time */
	pb := b.Bucket(sb("PRIVMSG"))
	if nil == pb {
		return 0, "", nil
	}
	last := pb.Get(sb("last"))
	if nil == last {
		return 0, "", nil
	}
	/* Make sure the user was active in the right time interval */
	if (-1 == bytes.Compare(last, from)) ||
		!(-1 == bytes.Compare(last, to)) {
		return 0, "", nil
	}
	/* Increase the bank account by the proper number of credits */
	cur, err := ChangeAccountBalance(b, PAYRATE)
	if nil != err {
		return 0, "", err
	}
	/* Get the last reply-to */
	what := pb.Get(sb("what"))
	if nil == what || 0 == len(what) {
		return 0, "", fmt.Errorf("%v said nothing, nowhere", v)
	}
	parts := strings.SplitN(string(what), " ", 2)
	if 2 != len(parts) {
		return 0, "", fmt.Errorf("not enough said")
	}

	log.Printf("[PAYROLL] Paid %v to %v (%s)", PAYRATE, v, last)
	return cur, parts[0], nil
}

/* CheckBalance sends a viewer's balance to him */
func CheckBalance(nick, replyto, args string) error {
	/* Get viewer's balance */
	b, err := GetAccountBalance(nick)
	go Privmsg(
		replyto,
		"%v: your balance is %v%v",
		nick,
		b,
		CURRENCYUNITS,
	)
	return err
}
