package main

/*
 * db.go
 * Database convenience functions
 * By MagisterQuis
 * Created 20160821
 * Last Modified 20160827
 */

import (
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
		bucket := GetBucket(tx, bs...)
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
func GetBool(k []byte, bs ...[]byte) *bool {
	var b *bool
	DB.Update(func(tx *bolt.Tx) error {
		b = GetBoolTx(tx, k, bs...)
		return nil
	})
	return b
}

/* GetBoolTx gets the bool in the key k in the bucket path bs in the
transaction tx.  It panics on error and returns nil if there was no value
stored. */
func GetBoolTx(tx *bolt.Tx, k []byte, bs ...[]byte) *bool {
	b := GetTx(tx, k, bs...)
	if nil == b {
		return nil
	}
	pb := byteBool(b)
	return &pb
}

/* PutTx puts a value v under key k in the bucket path named in bs in the
transaction Tx.  The buckets will be created if they do not exist. */
func PutTx(tx *bolt.Tx, k, v []byte, bs ...[]byte) {
	if err := GetBucket(tx, bs...).Put(k, v); nil != err {
		panic(err.Error())
	}

}

/* Get gets the value under key k stored in the buckets named in bs, which will
be created if they do not exist, all in the transaction Tx. */
func GetTx(tx *bolt.Tx, k []byte, bs ...[]byte) []byte {
	return GetBucket(tx, bs...).Get(k)
}

/* Put puts the value v in the key k in the buckets bs, which will be created
if they do not exist. */
func Put(k, v []byte, bs ...[]byte) {
	DB.Update(func(tx *bolt.Tx) error {
		PutTx(tx, k, v, bs...)
		return nil
	})
}

/* Get gets the value under k in the buckets bs, which will be created if they
do not exist. */
func Get(k []byte, bs ...[]byte) []byte {
	var b []byte
	DB.Update(func(tx *bolt.Tx) error {
		ib := GetTx(tx, k, bs...)
		if nil == ib {
			return nil
		}
		b = make([]byte, len(ib))
		copy(b, ib)
		return nil
	})
	return b
}

/* GetBucket gets the bucket with bucket path bs.  Buckets will be created if
they don't exist. */
func GetBucket(tx *bolt.Tx, bs ...[]byte) *bolt.Bucket {
	/* Make sure we actually have buckets */
	if nil == bs {
		panic("no bucket specified")
	}
	if 0 == len(bs) {
		panic("empty bucket path specified")
	}

	/* Get the initial bucket */
	bucket, err := tx.CreateBucketIfNotExists(bs[0])
	if nil != err {
		panic(err.Error())
	}

	/* Get subsequent buckets */
	for _, b := range bs[1:] {
		bucket, err = bucket.CreateBucketIfNotExists(b)
		if nil != err {
			panic(err.Error())
		}
	}
	return bucket
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
