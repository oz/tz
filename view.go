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

	"github.com/muesli/termenv"
)

// Width required to display 24 hours
const UIWidth = 94

func (m model) View() string {
	s := normalTextStyle("\n  What time is it?\n\n").String()

	// Show hours for each zone
	for zi, zone := range m.zones {
		hours := strings.Builder{}
		dates := strings.Builder{}

		startHour := 0
		if zi > 0 {
			startHour = (zone.currentTime().Hour() - m.zones[0].currentTime().Hour()) % 24
		}

		dateChanged := false
		for i := startHour; i < startHour+24; i++ {
			var hour int
			var period string

			// switch AM/PM vs 24-hour in future view
			if m.isTwelveHour {
				hour = ((i % 24) + 24) % 12
				if hour == 0 {
					hour = 12
				}
				if i%24 >= 12 {
					period = "PM"
				} else {
					period = "  "
				}
			} else {
				hour = ((i % 24) + 24) % 24 // mod 24
			}
			out := termenv.String(fmt.Sprintf("%2d%s", hour, period))

			out = out.Foreground(term.Color(hourColorCode(hour)))
			// Cursor
			if m.hour == i-startHour {
				out = out.Background(term.Color("#00B67F"))
				if hasDarkBackground {
					out = out.Foreground(term.Color("#262626")).Bold()
				} else {
					out = out.Foreground(term.Color("#f1f1f1"))
				}
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
			} else if hour == 12 && period != "PM" {
				dates.WriteString(formatDayChange(&m, zone))
				dateChanged = true
			}

			// handle the spacing for the date/calendar icon
			if !dateChanged && !m.isTwelveHour {
				dates.WriteString("    ")
			} else if m.isTwelveHour && !dateChanged {
				dates.WriteString("      ")
			}
		}

		// switch current timestamp from 12 to 24 hour
		var datetime string
		if m.isMilitary {
			datetime = zone.ShortMT()
		} else {
			datetime = zone.ShortDT()
		}

		zoneHeader := fmt.Sprintf("%s %s %s", zone.ClockEmoji(), normalTextStyle(zone.String()), dateTimeStyle(datetime))

		s += fmt.Sprintf("  %s\n  %s\n  %s\n", zoneHeader, hours.String(), dates.String())
	}

	if m.interactive {
		s += status()
	}
	return s
}

func status() string {
	text := "  q: quit, d: toggle date"
	for {
		text += " "
		if len(text) > UIWidth {
			text = text[0:UIWidth]
			break
		}
	}

	color := "#939183"
	if hasDarkBackground {
		color = "#605C5A"
	}

	status := termenv.String(text).Foreground(term.Color(color))

	return status.String()
}

func formatDayChange(m *model, z *Zone) string {
	zTime := z.currentTime()
	if zTime.Hour() > m.now.Hour() {
		zTime = zTime.AddDate(0, 0, 1)
	}

	color := "#777266"
	if hasDarkBackground {
		color = "#7B7573"
	}

	str := termenv.String(fmt.Sprintf("📆 %s", zTime.Format("Mon 02")))
	return str.Foreground(term.Color(color)).String()
}

// Return a color matching the time of the day at a given hour.
func hourColorCode(hour int) (color string) {
	switch hour {
	// Morning
	case 7, 8:
		if hasDarkBackground {
			color = "#98E1D8"
		} else {
			color = "#35B6A6"
		}

	// Day
	case 9, 10, 11, 12, 13, 14, 15, 16, 17:
		if hasDarkBackground {
			color = "#E8C64D"
		} else {
			color = "#FA8F2D"
		}

	// Evening
	case 18, 19:
		if hasDarkBackground {
			color = "#C95F48"
		} else {
			color = "#FC6442"
		}

	// Night
	default:
		if hasDarkBackground {
			color = "#5957C9"
		} else {
			color = "#664FC3"
		}
	}
	return color
}

func dateTimeStyle(str string) termenv.Style {
	color := "#777266"
	if hasDarkBackground {
		color = "#757575"
	}
	return termenv.String(str).Foreground(term.Color(color))
}

func normalTextStyle(str string) termenv.Style {
	var color = "#32312B"
	if hasDarkBackground {
		color = "#ECEAD9"
	}
	return termenv.String(str).Foreground(term.Color(color))
}
