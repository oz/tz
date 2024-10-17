/**
 * This file is part of tz.
 *
 * tz is free software: you can redistribute it and/or modify it under
 * the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your
 * option) any later version.
 *
 * tz is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
 * or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public
 * License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with tz.  If not, see <https://www.gnu.org/licenses/>.
 **/
package main

import (
	"testing"
	"time"
)

func TestNewClockNow(t *testing.T) {
	start := time.Now()
	clock := NewClockNow()
	if clock.t.Compare(start) < 0 {
		t.Errorf("Unexpected old time from NewClockNow(): %v", clock.t.Format(time.RFC3339))
	}
	if !clock.isRealTime {
		t.Error("NewClockNow() wasn’t real time")
	}
}

func TestClockTime(t *testing.T) {
	epoch, err := time.Parse(time.RFC3339, "1970-01-01T01:03:05Z")
	if err != nil {
		t.Fatalf("Could not parse test epoch")
	}
	clock := NewClockTime(epoch)
	if clock.t.Compare(epoch) != 0 {
		t.Errorf("Unexpected time from NewClockTime(epoch): %v", clock.t.Format(time.RFC3339))
	}
	if clock.isRealTime {
		t.Error("NewClockTime() shouldn’t be real time")
	}
}

func TestClockUnixTimestamp(t *testing.T) {
	epoch, err := time.Parse(time.RFC3339, "1970-01-01T00:00:01Z")
	if err != nil {
		t.Fatalf("Could not parse test epoch")
	}
	clock := NewClockUnixTimestamp(1)
	if clock.t.Compare(epoch) != 0 {
		t.Errorf("Unexpected time from NewClockUnixTimestamp(1): %v", clock.t.Format(time.RFC3339))
	}
	if clock.isRealTime {
		t.Error("NewClockUnixTimestamp() shouldn’t be real time")
	}
}
