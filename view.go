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
	"strings"
	"time"

	"github.com/muesli/termenv"
)

// Width required to display 24 hours
const UIWidth = 94

func (m model) View() string {
	s := "What time is it?\n\n"

	// Show hours for each zone
	for zi, zone := range m.zones {
		hours := strings.Builder{}
		dates := strings.Builder{}

		startHour := 0
		if zi > 0 {
			startHour = (zone.Offset - m.zones[0].Offset) % 24
		}

		dateChanged := false
		for i := startHour; i < startHour+24; i++ {
			hour := ((i % 24) + 24) % 24 // mod 24
			out := termenv.String(fmt.Sprintf("%2d", hour))

			out = out.Foreground(term.Color(hourColorCode(hour)))
			// Cursor
			if m.hour == i-startHour {
				out = out.Background(term.Color("41")).Foreground(term.Color("#000000"))
			}
			hours.WriteString(out.String())
			hours.WriteString("  ")

			// Show the day under the hour, when the date changes.
			if !m.showDates {
				continue
			}
			if hour == 0 {
				dates.WriteString(formatDayChange(&m, zone))
				dateChanged = true
			}
			if !dateChanged {
				dates.WriteString("    ")
			}
		}

		zoneHeader := termenv.String(fmt.Sprintf("%s %s: %s", zone.ClockEmoji(), zone, zone.ShortDT()))
		zoneHeader = zoneHeader.Background(term.Color("234")).Foreground(term.Color("255"))

		s += fmt.Sprintf("%s\n%s\n%s\n\n", zoneHeader, hours.String(), dates.String())
	}

	s += status()
	return s
}

func status() string {
	text := "q: quit, d: toggle date"
	for {
		text += " "
		if len(text) > UIWidth {
			text = text[0:UIWidth]
			break
		}
	}
	status := termenv.String(text).Background(term.Color("234")).Foreground(term.Color("255"))

	return status.String()
}

func formatDayChange(m *model, z *Zone) string {
	zTime := z.currentTime()
	if zTime.Hour() > time.Now().Hour() {
		zTime = zTime.AddDate(0, 0, 1)
	}
	str := termenv.String(fmt.Sprintf("📆%s", zTime.Format("Mon 02")))
	str = str.Foreground(term.Color("245"))
	return str.String()
}

// Return a color matching the time of the day at a given hour.
func hourColorCode(hour int) (color string) {
	switch hour {
	// Morning
	case 7, 8:
		color = "12"

	// Day
	case 9, 10, 11, 12, 13, 14, 15, 16, 17:
		color = "227"

	// Evening
	case 18, 19:
		color = "202"

	// Night
	default:
		color = "27"
	}
	return color
}
