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

type FormatStyle int

const (
	DefaultFormatStyle FormatStyle = iota
	IsoFormatStyle
	UnixFormatStyle
)

func (fs FormatStyle) next() FormatStyle {
	switch (fs) {
	case DefaultFormatStyle:
		return IsoFormatStyle
	case IsoFormatStyle:
		return UnixFormatStyle
	default:
		return DefaultFormatStyle
	}
}

func (fs FormatStyle) previous() FormatStyle {
	switch (fs) {
	case DefaultFormatStyle:
		return UnixFormatStyle
	case UnixFormatStyle:
		return IsoFormatStyle
	default:
		return DefaultFormatStyle
	}
}

type ZoneStyle int

const (
	AbbreviationZoneStyle ZoneStyle = iota
	WithZOffsetZoneStyle
	WithRelativeZoneStyle
)

func (zs ZoneStyle) next() ZoneStyle {
	switch (zs) {
	case AbbreviationZoneStyle:
		return WithZOffsetZoneStyle
	case WithZOffsetZoneStyle:
		return WithRelativeZoneStyle
	default:
		return AbbreviationZoneStyle
	}
}

func (zs ZoneStyle) previous() ZoneStyle {
	switch (zs) {
	case AbbreviationZoneStyle:
		return WithRelativeZoneStyle
	case WithRelativeZoneStyle:
		return WithZOffsetZoneStyle
	default:
		return AbbreviationZoneStyle
	}
}

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
	for i, zone := range m.zones {
		hours := strings.Builder{}
		dates := strings.Builder{}
		timeInZone := zone.currentTime(m.clock.t)
		midnightInZone := timeInZone.Add(-midnightOffset)
		wasDST := midnightInZone.Add(-time.Hour).IsDST()
		previousHour := midnightInZone.Add(-time.Hour).Hour()
		highlighted := i == (m.highlighted - 1)

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
		switch m.formatStyle {
		case IsoFormatStyle:
				datetime = timeInZone.Format("2006-01-02T15:04-07:00")
		case UnixFormatStyle:
			_, weekOfYear := timeInZone.ISOWeek()
			dayOfYear := timeInZone.Format("__2")
			yesNo := map[bool]string{true: "With", false: "No"}
			datetime = fmt.Sprintf(
				"%v DST, Week %v, Day %v, Unix %v",
				yesNo[timeInZone.IsDST()],
				weekOfYear,
				dayOfYear,
				timeInZone.Unix(),
			)
		default:
			if m.isMilitary {
				datetime = zone.ShortMT(m.clock.t)
			} else {
				datetime = zone.ShortDT(m.clock.t)
			}
		}

		var zoneString = zone.String(timeInZone)
		switch m.zoneStyle {
		case WithZOffsetZoneStyle:
			utcOffset := timeInZone.Format("Z-07:00")
			zoneString = fmt.Sprintf("[%s] %s", utcOffset, zoneString)
		case WithRelativeZoneStyle:
			_, otherOffset := timeInZone.Zone()
			_, localOffset := m.clock.t.Zone()
			relativeOffset := m.clock.t.In(time.FixedZone("", otherOffset - localOffset)).Format("-07:00")
			zoneString = fmt.Sprintf("[%s] %s", relativeOffset, zoneString)
		default:
		}

		clockString := zone.ClockEmoji(m.clock.t)
		usedZoneHeaderWidth := termenv.String(clockString + zoneString + datetime).Width()
		unusedZoneHeaderWidth := max(0, zoneHeaderWidth - usedZoneHeaderWidth - MinimumZoneHeaderPadding)
		rightAlignmentSpace := strings.Repeat(" ", unusedZoneHeaderWidth)
		zoneHeader := fmt.Sprintf("%s %s %s%s", clockString, normalTextStyle(zoneString), rightAlignmentSpace, dateTimeStyle(datetime))

		marker := "  "
		if highlighted {
			marker = termenv.String(">>").Reverse().String()
		}
		lines := []string{zoneHeader, hours.String(), dates.String()}
		for _, line := range lines {
			s += fmt.Sprintf("%s%s\n", marker, line)
		}
	}

	if m.interactive {
		s += status(m)
	}
	return s
}

func status(m model) string {

	var text []string

	if m.showHelp {
		text = []string{
			"?: help, -/+/0: minutes, h/l: hours, H/L: days, </>: weeks, t: go to now, j/k: highlight",
			"q: quit, d: toggle dates, f: toggle formats, z: toggle zone offsets, o: open in web",
		}
	} else {
		text = []string{
			"?: help",
			"q: quit",
		}
	}

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
