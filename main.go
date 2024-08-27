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
	"github.com/muesli/termenv"
)

// CurrentVersion represents the current build version.
const CurrentVersion = "0.7.0"

var (
	term              = termenv.ColorProfile()
	hasDarkBackground = termenv.HasDarkBackground()

	// Now is used around tz to share/set the current time.
	Now *Clock = NewClock(0)
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
	now         time.Time
	hour        int
	showDates   bool
	interactive bool
	isMilitary  bool
	watch       bool
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
			Now.AddHours(-1)

		case "right", "l":
			if m.hour > 22 {
				m.hour = 0
			} else {
				m.hour++
			}
			Now.AddHours(1)

		case "H":
			Now.AddDays(-1)

		case "L":
			Now.AddDays(1)

		case "<":
			Now.AddDays(-7)

		case ">":
			Now.AddDays(7)

		case "o":
			openInTimeAndDateDotCom(Now.Time())

		case "t":
			Now = NewClock(0)
			m.hour = Now.Time().Hour()

		case "d":
			m.showDates = !m.showDates
		}

	case tickMsg:
		if m.watch {
			m.now = time.Time(msg)
			Now = NewClock(0)
			m.hour = Now.Time().Hour()
		}
		return m, tick()
	}
	return m, nil
}

func main() {
	exitQuick := flag.Bool("q", false, "exit immediately")
	showVersion := flag.Bool("v", false, "show version")
	when := flag.Int64("when", 0, "time in seconds since unix epoch")
	doSearch := flag.Bool("list", false, "list zones by name")
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

	if *when != 0 {
		Now = NewClock(*when)
	}
	config, err := LoadConfig(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %s\n", err)
		os.Exit(2)
	}
	var initialModel = model{
		zones:      config.Zones,
		now:        Now.Time(),
		hour:       Now.Time().Hour(),
		showDates:  false,
		isMilitary: *military,
		watch:      *watch,
	}

	initialModel.interactive = !*exitQuick

	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
