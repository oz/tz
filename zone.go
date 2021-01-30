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

import "time"

var name, offset = time.Now().Zone()
var DefaultZones = []Zone{
	{
		name:   "Local",
		dbName: name,
		offset: offset / 3600,
	},
	{
		name:   "Paris",
		dbName: "Europe/Paris",
		offset: 1,
	},
	{
		name:   "New-York",
		dbName: "America/New_York",
		offset: -5,
	},
	{
		name:   "UTC",
		dbName: "UTC",
		offset: 0,
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
	dbName string // tz db name
	name   string // custom name
	offset int    // Integer offset from UTC
}

func (z Zone) String() string {
	return z.name
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
	now := time.Now()
	zName, _ := now.Zone()
	if z.dbName != zName {
		loc, err := time.LoadLocation(z.dbName)
		if err != nil {
			return now
		}
		return now.In(loc)
	}
	return now
}
