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
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var term = termenv.ColorProfile()

type model struct {
	zones []Zone
	hour  int
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "left", "h":
			if m.hour == 0 {
				m.hour = 23
			} else {
				m.hour--
			}

		case "right", "l":
			if m.hour > 22 {
				m.hour = 0
			} else {
				m.hour++
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "What time is it?\n\n"

	// Show hours for each zone
	for zi, z := range m.zones {
		hours := strings.Builder{}
		startAt := 0
		if zi > 0 {
			startAt = (z.offset - m.zones[0].offset) % 24
		}

		// A list of hours
		for i := startAt; i < startAt+24; i++ {
			hour := ((i % 24) + 24) % 24
			out := termenv.String(fmt.Sprintf("%2d", hour))

			out = out.Foreground(term.Color(hourColorCode(hour)))
			// Cursor
			if m.hour == i-startAt {
				out = out.Background(term.Color("41")).Foreground(term.Color("#000000"))
			}
			hours.WriteString(out.String())
			hours.WriteString("  ")
		}

		zoneHeader := termenv.String(fmt.Sprintf("%s %s: %s", z.ClockEmoji(), z, z.ShortDT()))
		zoneHeader = zoneHeader.Background(term.Color("234")).Foreground(term.Color("255"))

		s += fmt.Sprintf("%s\n%s\n\n", zoneHeader, hours.String())
	}

	s += "\nPress q to quit.\n"
	return s
}

// Return a color matching the time of the day at a given hour.
func hourColorCode(hour int) (color string) {
	switch hour {
	// Morning
	case 7, 8:
		color = "12"

	// Day
	case 9, 10, 11, 12, 13, 14, 15, 16, 17:
		color = "11"

	// Evening
	case 18, 19:
		color = "3"

	// Night
	default:
		color = "17"
	}
	return color
}

func main() {
	now := time.Now()
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s", err)
		os.Exit(2)
	}
	var initialModel = model{
		zones: config.Zones,
		hour:  now.Hour(),
	}
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
