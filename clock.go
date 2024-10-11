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

// AddDays adds n hours to the current date-time, keeping the minutes
func (c *Clock) AddHours(n int) {
	c.t = c.t.Add(time.Hour * time.Duration(n))
	c.isRealTime = false
}

// AddMinutes adds n minutes to the current date-time, keeping the seconds
func (c *Clock) AddMinutes(n int) {
	c.t = c.t.Add(time.Minute * time.Duration(n))
	c.isRealTime = false
}

// Get the wrapped time.Time struct
func (c *Clock) Time() time.Time {
	return c.t
}
