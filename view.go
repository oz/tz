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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/muesli/termenv"
	xterm "golang.org/x/term"
)

// Width required to display 24 hours
const UIWidth = 94
const MinimumZoneHeaderPadding = 6
const MaximumZoneHeaderColumns = UIWidth + MinimumZoneHeaderPadding

func (m model) View() string {
	s := normalTextStyle("\n  What time is it?\n\n").String()

	zoneHeaderWidth := MaximumZoneHeaderColumns
	envWidth, envErr := strconv.Atoi(os.Getenv("COLUMNS"))
	if envErr == nil {
		zoneHeaderWidth = min(envWidth, zoneHeaderWidth)
	} else {
		fd := int(os.Stdout.Fd())
		if xterm.IsTerminal(fd) {
			termWidth, _, termErr := xterm.GetSize(fd)
			if termErr == nil {
				zoneHeaderWidth = min(termWidth, zoneHeaderWidth)
			}
		}
	}

	midnight := time.Date(
		m.clock.t.Year(),
		m.clock.t.Month(),
		m.clock.t.Day(),
		0, // Hours
		m.clock.t.Minute(),
		0, // Seconds
		0, // Nanoseconds
		m.clock.t.Location(),
	)
	midnightOffset := time.Duration(m.clock.t.UnixNano() - midnight.UnixNano())
	cursorColumn := int(midnightOffset / time.Hour)

	// Show hours for each zone
	for _, zone := range m.zones {
		hours := strings.Builder{}
		dates := strings.Builder{}
		timeInZone := zone.currentTime(m.clock.t)
		midnightInZone := timeInZone.Add(-midnightOffset)
		wasDST := midnightInZone.Add(-time.Hour).IsDST()
		previousHour := midnightInZone.Add(-time.Hour).Hour()

		dateChanged := false
		for column := 0; column < 24; column++ {
			time := midnightInZone.Add(time.Duration(column) * time.Hour)
			nowDST := time.IsDST()
			hour := time.Hour()
			out := termenv.String(fmt.Sprintf("%2d", hour))

			out = out.Foreground(term.Color(hourColorCode(hour)))
			// Cursor
			if column == cursorColumn {
				out = out.Background(term.Color(hourColorCode(hour)))
				if hasDarkBackground {
					out = out.Foreground(term.Color("#262626")).Bold()
				} else {
					out = out.Foreground(term.Color("#f1f1f1"))
				}
			}
			hours.WriteString(out.String())
			hours.WriteString("  ")

			// Show the day under the hour, when the date changes.
			if m.showDates {
				if hour < previousHour {
					dates.WriteString(formatDayChange(&m, zone))
					dateChanged = true
				}

				if wasDST != nowDST {
					if nowDST {
						dates.WriteString("=DST")
					} else {
						dates.WriteString("≠DST")
					}
				} else if !dateChanged {
					dates.WriteString("    ")
				}
			}

			wasDST = nowDST
			previousHour = hour
		}

		var datetime string
		if m.isMilitary {
			datetime = zone.ShortMT(m.clock.t)
		} else {
			datetime = zone.ShortDT(m.clock.t)
		}

		clockString := zone.ClockEmoji(m.clock.t)
		zoneString := zone.VerboseString(m.clock.t)
		usedZoneHeaderWidth := termenv.String(clockString + zoneString + datetime).Width()
		unusedZoneHeaderWidth := max(0, zoneHeaderWidth - usedZoneHeaderWidth - MinimumZoneHeaderPadding)
		rightAlignmentSpace := strings.Repeat(" ", unusedZoneHeaderWidth)
		zoneHeader := fmt.Sprintf("%s %s %s%s", clockString, normalTextStyle(zoneString), rightAlignmentSpace, dateTimeStyle(datetime))

		s += fmt.Sprintf("  %s\n  %s\n  %s\n", zoneHeader, hours.String(), dates.String())
	}

	if m.interactive {
		s += status(m)
	}
	return s
}

// Generate the help lines
func generateKeymapStrings(k Keymaps, showAll bool) []string {
	helpKey := fmt.Sprintf("%s: help", k.Help[0])
	quitKey := fmt.Sprintf("%s: quit", k.Quit[0])

	if showAll {
		delimiter := ", "
		return []string {
			strings.Join(
				[]string {
					helpKey,
					fmt.Sprintf("%s/%s/%s: minutes", k.PrevMinute[0], k.NextMinute[0], k.ZeroMinute[0]),
					fmt.Sprintf("%s/%s: hours", k.PrevHour[0], k.NextHour[0]),
					fmt.Sprintf("%s/%s: days", k.PrevDay[0], k.NextDay[0]),
					fmt.Sprintf("%s/%s: weeks", k.PrevWeek[0], k.NextWeek[0]),
					fmt.Sprintf("%s: go to now", k.Now[0]),
				},
				delimiter,
			),
			strings.Join(
				[]string {
					quitKey,
					fmt.Sprintf("%s: toggle dates", k.ToggleDate[0]),
					fmt.Sprintf("%s: open in web", k.OpenWeb[0]),
				},
				delimiter,
			),
		}
	} else {
		return []string {
			helpKey,
			quitKey,
		}
	}
}

func status(m model) string {
	var text []string = generateKeymapStrings(m.keymaps, m.showHelp)

	backgroundPadding := strings.Repeat(" ", UIWidth)
	for i, line := range text {
		text[i] = ("  " + line + backgroundPadding)[:UIWidth]
	}

	color := "#939183"
	if hasDarkBackground {
		color = "#605C5A"
	}

	status := termenv.String(strings.Join(text, "\n")).Foreground(term.Color(color))

	return status.String()
}

func formatDayChange(m *model, z *Zone) string {
	zTime := z.currentTime(m.clock.t)
	if zTime.Hour() > m.clock.t.Hour() {
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
