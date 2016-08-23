package main

/*
 * db.go
 * Database convenience functions
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160821
 */

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/boltdb/bolt"
)

/* Trap ^C and write database before we go */
func CatchInt() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Printf("Caught interrupt, closing database")
		if err := DB.Close(); nil != err {
			log.Fatalf("Unable to close database: %v", err)
		}
		log.Fatalf("Closed database.  Bye.")
	}()
}

/* PutBool stores the v in k in b */
func PutBool(b *bolt.Bucket, k []byte, v bool) error {
	return b.Put(k, []byte{0x01})
}

/* GetBool gets the bool in k in b.  If the value wasn't found, ok will be
false. */
func GetBool(b *bolt.Bucket, k []byte) (value, ok bool, err error) {
	/* Get the value from the database */
	v := b.Get(k)
	if nil == v {
		/* Not in there */
		return false, false, nil
	}
	/* Make sure it's what we expect */
	if 1 != len(v) {
		return false, false, fmt.Errorf(
			"value was incorrect length (%v)",
			len(v),
		)
	}
	/* Booleanize it */
	switch v[0] {
	case 0x01:
		return true, true, nil
	case 0x00:
		return false, true, nil
	default:
		return false, false, fmt.Errorf(
			"value has unexpected value 0x(%02X)",
			v[0],
		)
	}
}

/* bs turns a string into a byte slice */
func bs(b []byte) string { return string(b) }

/* sb turns a byte slice into a string */
func sb(s string) []byte { return []byte(s) }
