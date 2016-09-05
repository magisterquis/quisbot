package main

/*
 * ratelimit.go
 * Rate-limiter for long messages
 * By J. Stuart McMurray
 * Created 20160901
 * Last Modified 20160901
 */

import (
	"sync"
	"time"
)

/* A RateLimiter is used to keep long messages from being sent too often. */
type RateLimiter struct {
	next time.Time
	wait time.Duration
	lock *sync.Mutex
}

/* NewRateLimiter returns a new RateLimiter which won't let messages be sent
within wait time of each other. */
func NewRateLimiter(wait time.Duration) *RateLimiter {
	return &RateLimiter{
		wait: wait,
		lock: &sync.Mutex{},
	}
}

/* Until returns the amount of time until the next message can be sent. */
func (r *RateLimiter) Until() time.Duration {
	r.lock.Lock()
	defer r.lock.Unlock()
	now := time.Now()
	/* If we have no wait, note the next time */
	if now.After(r.next) {
		r.next = now.Add(r.wait)
		return 0
	}
	return r.next.Sub(now)
}
