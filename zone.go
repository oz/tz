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
	"fmt"
	"time"
)

var name, _ = time.Now().Zone()
var DefaultZones = []*Zone{
	{
		Name:   "Local",
		DbName: name,
	},
	{
		Name:   "UTC",
		DbName: "UTC",
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

// Zone stores the name of a time zone
type Zone struct {
	Loc    *time.Location
	DbName string // Name in tzdata
	Name   string // Preferred name (user-provided, or else DbName by default)
}

func (z Zone) String(t time.Time) string {
	return fmt.Sprintf("(%s) %s", z.Abbreviation(t), z.Name)
}

// Abbreviated short name for the zone (e.g. acronym "ABC" if available, or else a number like "-3").
// It depends dynamically on the daylight saving policy in the zone at time `t`.
func (z Zone) Abbreviation(t time.Time) string {
	shortName, _ := z.currentTime(t).Zone()
	return shortName
}

// ClockEmoji returns the corresponding emoji clock for a given hour
func (z Zone) ClockEmoji(t time.Time) string {
	h := ((z.currentTime(t).Hour() % 12) + 12) % 12
	return EmojiClocks[h]
}

// ShortDT returns the current time in short format.
func (z Zone) ShortDT(t time.Time) string {
	return z.currentTime(t).Format("3:04PM, Mon Jan 02, 2006")
}

// ShortMT returns the current military time in short format.
func (z Zone) ShortMT(t time.Time) string {
	return z.currentTime(t).Format("15:04, Mon Jan 02, 2006")
}

func (z Zone) currentTime(t time.Time) time.Time {
	zName, _ := t.Zone()
	if z.DbName != zName {
		return t.In(z.Loc)
	}
	return t
}
