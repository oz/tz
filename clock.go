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

// AddDays adds n days to the current date and clears the minutes
func (c *Clock) AddHours(n int) {
	c.t = time.Date(
		c.t.Year(),
		c.t.Month(),
		c.t.Day(),
		c.t.Hour(),
		0, // Minutes set to 0
		0, // Seconds set to 0
		0, // Nanoseconds set to 0
		c.t.Location(),
	)
	c.t = c.t.Add(time.Hour * time.Duration(n))
}

// Get the wrapped time.Time struct
func (c *Clock) Time() time.Time {
	return c.t
}
