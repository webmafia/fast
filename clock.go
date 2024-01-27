package fast

import (
	"sync/atomic"
	"time"
)

type Clock struct {
	ts     int64
	closer chan struct{}
}

func NewClock() *Clock {
	c := &Clock{
		ts:     time.Now().Unix(),
		closer: make(chan struct{}),
	}

	go c.tick()

	return c
}

//go:inline
func (c *Clock) Now() time.Time {
	return time.Unix(c.Unix(), 0)
}

//go:inline
func (c *Clock) Unix() int64 {
	return atomic.LoadInt64(&c.ts)
}

func (c *Clock) Close() {
	c.closer <- struct{}{}
}

func (c *Clock) tick() {
	ticker := time.NewTicker(time.Second)
	now := time.Now()

	// Wait to the next second
	nextSec := now.Truncate(time.Second).Add(time.Second)
	delay := nextSec.Sub(now)

	if delay > 0 {
		ticker.Reset(delay)
		<-ticker.C
		ticker.Reset(time.Second)
	}

	for {
		select {
		case <-c.closer:
			return
		case ts := <-ticker.C:
			atomic.StoreInt64(&c.ts, ts.Unix())
		}
	}
}
