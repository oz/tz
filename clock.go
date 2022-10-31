package main

import "time"

// Clock keeps track of the current time.
type Clock struct {
	t time.Time
}

// Make a new Clock.
//
// When "sec" is not t0, the returned Clock is set to this time,
// assuming "sec" is a UNIX timestamp (cf. time.Unix).
func NewClock(sec int64) *Clock {
	clock := new(Clock)

	if sec != 0 {
		clock.t = time.Unix(sec, 0)
	} else {
		clock.t = time.Now()
	}
	return clock
}

// AddDays adds n days to the current date.
func (c *Clock) AddDays(n int) {
	c.t = c.t.AddDate(0, 0, n)
}

// Get the wrapped time.Time struct
func (c *Clock) Time() time.Time {
	return c.t
}
