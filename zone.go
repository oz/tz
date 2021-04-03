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
	"strconv"
	"time"

	"github.com/WIZARDISHUNGRY/tzcoords"
	"github.com/kelvins/sunrisesunset"
)

var name, offset = time.Now().Zone()
var DefaultZones = []*Zone{
	{
		Name:   "Local",
		DbName: name,
		Offset: offset / 3600,
	},
	{
		Name:   "UTC",
		DbName: "UTC",
		Offset: 0,
	},
}

// EmojiClocks maps hour to the corresponding clock emoji
var EmojiClocks = map[int]string{
	0:  "🕛",
	1:  "🕐",
	2:  "🕑",
	3:  "🕒",
	4:  "🕓",
	5:  "🕔",
	6:  "🕕",
	7:  "🕖",
	8:  "🕗",
	9:  "🕘",
	10: "🕙",
	11: "🕙",
}

// Zone stores the name of a time zone and its integer offset from UTC.
type Zone struct {
	DbName string // Name in tzdata
	Name   string // Short name
	Offset int    // Integer offset from UTC, in hours.
}

func (z Zone) String() string {
	return z.Name
}

// ClockEmoji returns the corresponding emoji clock for a given hour
func (z Zone) ClockEmoji() (clock string) {
	h := ((z.currentTime().Hour() % 12) + 12) % 12
	return EmojiClocks[h]
}

// ShortDT returns the current time in short format.
func (z Zone) ShortDT() string {
	return z.currentTime().Format("3:04PM, Mon 02")
}

func (z Zone) currentTime() time.Time {
	now := clock()
	zName, _ := now.Zone()
	if z.DbName != zName {
		loc, err := time.LoadLocation(z.DbName)
		if err != nil {
			return now
		}
		return now.In(loc)
	}
	return now
}

// SunriseSunset returns the sunrise and sunset for a zone
func (z Zone) SunriseSunset() (time.Time, time.Time, error) {
	lat, lon, err := tzcoords.ByString(z.DbName)
	now := z.currentTime()
	if err != nil {
		return now, now, err
	}
	utc, err := strconv.Atoi(now.Format("-0700"))
	if err != nil {
		return now, now, err
	}
	utcFloat := float64(utc) / 100.0
	p := sunrisesunset.Parameters{
		Latitude:  lat,
		Longitude: lon,
		UtcOffset: utcFloat,
		Date:      now,
	}
	return p.GetSunriseSunset()
}

// LightCycle returns a color schedule
func (z Zone) LightCycle() LightCycle {
	rise, set, err := z.SunriseSunset()
	if err != nil {
		return defaultHourHashes()
	}
	c := LightCycle{
		morning: make(map[int]struct{}),
		day:     make(map[int]struct{}),
		evening: make(map[int]struct{}),
		night:   make(map[int]struct{}),
	}
	var next map[int]struct{}
	for hour := 0; hour < 24; hour++ {
		if hour < rise.Hour() {
			next = c.night
		} else if hour == rise.Hour() {
			next = c.morning
		} else if hour > rise.Hour() {
			if hour < set.Hour() {
				next = c.day
			} else if hour == set.Hour() {
				next = c.evening
			} else {
				next = c.night
			}
		}
		next[hour] = struct{}{}
	}
	return c
}

// LightCycle defines the different periods of the day
type LightCycle struct {
	morning, day, evening, night map[int]struct{}
}

func defaultHourHashes() LightCycle {
	c := LightCycle{
		morning: make(map[int]struct{}),
		day:     make(map[int]struct{}),
		evening: make(map[int]struct{}),
		night:   make(map[int]struct{}),
	}
	for hour := 0; hour < 24; hour++ {
		switch hour {
		// Morning
		case 7, 8:
			c.morning[hour] = struct{}{}

		// Day
		case 9, 10, 11, 12, 13, 14, 15, 16, 17:
			c.day[hour] = struct{}{}

		// Evening
		case 18, 19:
			c.evening[hour] = struct{}{}

		// Night
		default:
			c.night[hour] = struct{}{}

		}
	}
	return c
}
