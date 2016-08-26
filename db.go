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

/* PutTx puts a value v under key k in the bucket path named in bs in the
transaction Tx.  The buckets will be created if they do not exist. */
func PutTx(tx *bolt.Tx, bs [][]byte, k, v []byte) error {
	b, err := GetBucket(tx, bs)
	if nil != err {
		return err
	}
	return b.Put(k, v)

}

/* Get gets the value under key k stored in the buckets named in bs, which will
be created if they do not exist, all in the transaction Tx. */
func GetTx(tx *bolt.Tx, bs [][]byte, k []byte) ([]byte, error) {
	b, err := GetBucket(tx, bs)
	if nil != err {
		return nil, err
	}
	return b.Get(k), nil
}

/* Put puts the value v in the key k in the buckets bs, which will be created
if they do not exist. */
func Put(bs [][]byte, k, v []byte) error {
	return DB.Update(func(tx *bolt.Tx) error {
		return PutTx(tx, bs, k, v)
	})
}

/* Get gets the value under k in the buckets bs, which will be created if they
do not exist. */
func Get(bs [][]byte, k []byte) ([]byte, error) {
	var b []byte
	derr := DB.Update(func(tx *bolt.Tx) error {
		ib, err := GetTx(tx, bs, k)
		if nil != err {
			return err
		}
		b = make([]byte, len(ib))
		copy(b, ib)
		return nil
	})
	return b, derr
}

/* GetBucket gets the bucket with bucket path bs */
func GetBucket(tx *bolt.Tx, bs [][]byte) (*bolt.Bucket, error) {
	/* Make sure we actually have buckets */
	if nil == bs || 0 == len(bs) {
		return nil, fmt.Errorf("no buckets specified")
	}

	/* Get the initial bucket */
	bucket, err := tx.CreateBucketIfNotExists(bs[0])
	if nil != err {
		return nil, err
	}

	/* Get subsequent buckets */
	for _, b := range bs[1:] {
		bucket, err = bucket.CreateBucketIfNotExists(b)
		if nil != err {
			return nil, err
		}
	}
	return bucket, nil
}

/* bs turns a string into a byte slice */
func bs(b []byte) string { return string(b) }

/* sb turns a byte slice into a string */
func sb(s string) []byte { return []byte(s) }
