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
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var (
	term              = termenv.ColorProfile()
	hasDarkBackground = termenv.HasDarkBackground()
)

type tickMsg time.Time

// Send a tickMsg every minute, on the minute.
func tick() tea.Cmd {
	return tea.Every(time.Minute, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	zones       []*Zone
	now         time.Time
	hour        int
	showDates   bool
	interactive bool
	time24      bool
}

func (m model) Init() tea.Cmd {
	// If -q flag is passed, send quit message after first render.
	if !m.interactive {
		return tea.Quit
	}

	// Fire initial tick command to begin receiving ticks on the minute.
	return tick()
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

		case "d":
			m.showDates = !m.showDates
		}

	case tickMsg:
		m.now = time.Time(msg)
		return m, tick()
	}
	return m, nil
}

func main() {
	now := time.Now()
	config, err := LoadConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Config error: %s\n", err)
		os.Exit(2)
	}
	var initialModel = model{
		zones:     config.Zones,
		now:       now,
		hour:      now.Hour(),
		showDates: false,
		time24:    config.Time24,
	}
	initialModel.interactive = !config.ExitQuick

	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
