package main

/*
 * db.go
 * Database convenience functions
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160826
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

/* PutBool stores v in the key k in the bucket path bs, and returns the
previous value, or nil if there was no previous value. */
func PutBool(k []byte, v bool, bs ...[]byte) (*bool, error) {
	var prev *bool
	derr := DB.Update(func(tx *bolt.Tx) error {
		bucket, err := GetBucket(tx, bs...)
		if nil != err {
			return err
		}
		/* Get previous value */
		p := bucket.Get(k)
		if nil != p {
			b := byteBool(p)
			prev = &b
		}
		/* Stick in the current value */
		return bucket.Put(k, boolByte(v))
	})
	return prev, derr
}

/* GetBool gets the bool in the key k in the bucket path bs.  It panics if v
is not a proper bool.  It returns nil if there was no value stored. */
func GetBool(k []byte, bs ...[]byte) (*bool, error) {
	b, err := Get(k, bs...)
	if nil != err {
		return nil, err
	}
	if nil == b {
		return nil, nil
	}
	pb := byteBool(b)
	return &pb, nil
}

/* PutTx puts a value v under key k in the bucket path named in bs in the
transaction Tx.  The buckets will be created if they do not exist. */
func PutTx(tx *bolt.Tx, k, v []byte, bs ...[]byte) error {
	b, err := GetBucket(tx, bs...)
	if nil != err {
		return err
	}
	return b.Put(k, v)

}

/* Get gets the value under key k stored in the buckets named in bs, which will
be created if they do not exist, all in the transaction Tx. */
func GetTx(tx *bolt.Tx, k []byte, bs ...[]byte) ([]byte, error) {
	b, err := GetBucket(tx, bs...)
	if nil != err {
		return nil, err
	}
	return b.Get(k), nil
}

/* Put puts the value v in the key k in the buckets bs, which will be created
if they do not exist. */
func Put(k, v []byte, bs ...[]byte) error {
	return DB.Update(func(tx *bolt.Tx) error {
		return PutTx(tx, k, v, bs...)
	})
}

/* Get gets the value under k in the buckets bs, which will be created if they
do not exist. */
func Get(k []byte, bs ...[]byte) ([]byte, error) {
	var b []byte
	derr := DB.Update(func(tx *bolt.Tx) error {
		ib, err := GetTx(tx, k, bs...)
		if nil != err {
			return err
		}
		if nil == ib {
			return nil
		}
		b = make([]byte, len(ib))
		copy(b, ib)
		return nil
	})
	return b, derr
}

/* GetBucket gets the bucket with bucket path bs */
func GetBucket(tx *bolt.Tx, bs ...[]byte) (*bolt.Bucket, error) {
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

/* byteBool turns a byte slice into a bool.  It panics if the byte slice
contains anything besides a single 0 or 1. */
func byteBool(b []byte) bool {
	if 1 != len(b) {
		log.Panicf("Improperly-sized bool []byte: %q", b)
	}
	if 0x00 == b[0] {
		return false
	}
	if 0x01 == b[0] {
		return true
	}
	log.Panicf("Improper bool []byte: %q", b)
	/* Can't get here */
	return false
}

/* boolByte returns a byte slice represeting b */
func boolByte(b bool) []byte {
	if b {
		return []byte{0x01}
	}
	return []byte{0x00}
}

/* bs turns a string into a byte slice */
func bs(b []byte) string { return string(b) }

/* sb turns a byte slice into a string */
func sb(s string) []byte { return []byte(s) }
