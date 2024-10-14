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
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
)

// CurrentVersion represents the current build version.
const CurrentVersion = "0.7.0"

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

func openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

func openInTimeAndDateDotCom(t time.Time) error {
	utcTime := t.In(time.UTC).Format("20060102T150405")
	url := fmt.Sprintf("https://www.timeanddate.com/worldclock/converter.html?iso=%s&p1=1440", utcTime)

	return openURL(url)
}

type model struct {
	zones       []*Zone
	clock       Clock
	highlighted int // 0 == none, else row number indexed from 1
	showDates   bool
	interactive bool
	isMilitary  bool
	watch       bool
	showHelp    bool
	formatStyle FormatStyle
	zoneStyle   ZoneStyle
}

func (m model) Init() tea.Cmd {
	// If -q flag is passed, send quit message after first render.
	if !m.interactive {
		return tea.Quit
	}

	// Fire initial tick command to begin receiving ticks on the minute.
	return tick()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "-":
			m.clock.AddMinutes(-1)

		case "+":
			m.clock.AddMinutes(1)

		case "0":
			m.clock = *NewClockTime(time.Date(
				m.clock.t.Year(),
				m.clock.t.Month(),
				m.clock.t.Day(),
				m.clock.t.Hour(),
				0,
				0,
				0,
				m.clock.t.Location(),
			))

		case "up", "k":
			modulo := len(m.zones) + 1
			m.highlighted = (m.highlighted - 1 + modulo) % modulo

		case "down", "j":
			modulo := len(m.zones) + 1
			m.highlighted = (m.highlighted + 1) % modulo

		case "left", "h":
			m.clock.AddHours(-1)

		case "right", "l":
			m.clock.AddHours(1)

		case "H":
			m.clock.AddDays(-1)

		case "L":
			m.clock.AddDays(1)

		case "<":
			m.clock.AddDays(-7)

		case ">":
			m.clock.AddDays(7)

		case "f":
			m.formatStyle = m.formatStyle.next()

		case "F":
			m.formatStyle = m.formatStyle.previous()

		case "o":
			openInTimeAndDateDotCom(m.clock.Time())

		case "t":
			m.clock = *NewClockNow()

		case "?":
			m.showHelp = !m.showHelp

		case "d":
			m.showDates = !m.showDates

		case "z":
			m.zoneStyle = m.zoneStyle.next()

		case "Z":
			m.zoneStyle = m.zoneStyle.previous()
		}

	case tickMsg:
		if m.watch && m.clock.isRealTime {
			m.clock = *NewClockNow()
		}
		return m, tick()
	}
	return m, nil
}

func main() {
	exitQuick := flag.Bool("q", false, "exit immediately")
	showVersion := flag.Bool("v", false, "show version")
	when := flag.Int64("when", 0, "time in seconds since unix epoch (disables -w)")
	doSearch := flag.Bool("list [filter]", false, "list zones by name")
	military := flag.Bool("m", false, "use 24-hour time")
	watch := flag.Bool("w", false, "watch live, set time to now every minute")
	flag.Parse()

	if *showVersion == true {
		fmt.Printf("tz %s\n", CurrentVersion)
		os.Exit(0)
	}

	if *doSearch {
		q := ""
		if arg := flag.Arg(0); arg != "" {
			q = arg
		}
		results := SearchZones(strings.ToLower(q))
		results.Print(os.Stdout)
		os.Exit(0)
	}

	config, err := LoadConfig(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %s\n", err)
		os.Exit(2)
	}
	var initialModel = model{
		zones:      config.Zones,
		clock:      *NewClockNow(),
		showDates:  false,
		isMilitary: *military,
		watch:      *watch,
		showHelp:   false,
		zoneStyle:  AbbreviationZoneStyle,
	}

	if *when != 0 {
		initialModel.clock = *NewClockUnixTimestamp(*when)
	}

	initialModel.interactive = !*exitQuick && isatty.IsTerminal(os.Stdout.Fd())

	p := tea.NewProgram(&initialModel)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
