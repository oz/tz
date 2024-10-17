package main

import "time"

// Clock keeps track of the current time.
type Clock struct {
	t time.Time
	isRealTime bool
}

// A new Clock in the local timezone at the current time now.
func NewClockNow() *Clock {
	clock := new(Clock)
	clock.t = time.Now()
	clock.isRealTime = true
	return clock
}

// A new Clock with the given time and timezone.
func NewClockTime(t time.Time) *Clock {
	clock := new(Clock)
	clock.t = t
	clock.isRealTime = false
	return clock
}

// A new Clock in the local timezone at the `time.Unix` timestamp.
func NewClockUnixTimestamp(sec int64) *Clock {
	clock := new(Clock)
	clock.t = time.Unix(sec, 0)
	clock.isRealTime = false
	return clock
}

// AddDays adds n days to the current date.
func (c *Clock) AddDays(n int) {
	c.t = c.t.AddDate(0, 0, n)
	c.isRealTime = false
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
	c.isRealTime = false
}

// Get the wrapped time.Time struct
func (c *Clock) Time() time.Time {
	return c.t
}
