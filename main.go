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
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
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

type whenValue struct {
	when **int64
}

func (w whenValue) String() string {
	if w.when != nil && *w.when != nil {
		return strconv.FormatInt(**w.when, 10)
	}
	return ""
}

func (w whenValue) Set(s string) error {
	if i, err := strconv.ParseInt(s, 0, 64); err == nil {
		*w.when = &i
		return nil
	} else if t, err := time.Parse(time.RFC3339, s); err == nil {
		u := t.Unix()
		*w.when = &u
		return nil
	}
	return errors.New("Could not parse integer (format: 123) or date-time (format: 2006-01-02T15:04:05+07:00)")
}

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
	keymaps     Keymaps
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

func match(input string, options []string) bool {
	return slices.Contains(options, input)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		key := msg.String()
		switch {

		case match(key, m.keymaps.Quit):
			return m, tea.Quit

		case match(key, m.keymaps.PrevMinute):
			m.clock.AddMinutes(-1)

		case match(key, m.keymaps.NextMinute):
			m.clock.AddMinutes(1)

		case match(key, m.keymaps.ZeroMinute):
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

		case match(key, m.keymaps.PrevHour):
			m.clock.AddHours(-1)

		case match(key, m.keymaps.NextHour):
			m.clock.AddHours(1)

		case match(key, m.keymaps.PrevDay):
			m.clock.AddDays(-1)

		case match(key, m.keymaps.NextDay):
			m.clock.AddDays(1)

		case match(key, m.keymaps.PrevWeek):
			m.clock.AddDays(-7)

		case match(key, m.keymaps.NextWeek):
			m.clock.AddDays(7)

		case match(key, m.keymaps.PrevLine):
			modulo := len(m.zones) + 1
			m.highlighted = (m.highlighted - 1 + modulo) % modulo

		case match(key, m.keymaps.NextLine):
			modulo := len(m.zones) + 1
			m.highlighted = (m.highlighted + 1) % modulo

		case match(key, m.keymaps.NextFStyle):
			m.formatStyle = m.formatStyle.next()

		case match(key, m.keymaps.PrevFStyle):
			m.formatStyle = m.formatStyle.previous()

		case match(key, m.keymaps.PrevZStyle):
			m.zoneStyle = m.zoneStyle.previous()

		case match(key, m.keymaps.NextZStyle):
			m.zoneStyle = m.zoneStyle.next()

		case match(key, m.keymaps.OpenWeb):
			openInTimeAndDateDotCom(m.clock.Time())

		case match(key, m.keymaps.Now):
			m.clock = *NewClockNow()

		case match(key, m.keymaps.ToggleDate):
			m.showDates = !m.showDates

		case match(key, m.keymaps.Help):
			m.showHelp = !m.showHelp
		}

	case tickMsg:
		if m.watch && m.clock.isRealTime {
			m.clock = *NewClockNow()
		}
		return m, tick()
	}
	return m, nil
}

func parseMainArgs() *model {
	logger.Println("Startup")

	var when *int64
	flag.Var(whenValue{&when}, "when", "date-time in seconds since unix epoch, or in ISO8601/RFC3339 format (disables -w)")
	exitQuick := flag.Bool("q", false, "exit immediately")
	showVersion := flag.Bool("v", false, "show version")
	doSearch := flag.Bool("list", false, "[filter] list or search zones by name")
	military := flag.Bool("m", false, "use 24-hour time")
	watch := flag.Bool("w", false, "watch live, set time to now every minute")
	flag.Parse()

	if *showVersion {
		fmt.Fprintf(os.Stdout(), "tz %s\n", CurrentVersion)
		os.Exit(0)
		return nil
	}

	if *doSearch {
		q := ""
		if arg := flag.Arg(0); arg != "" {
			q = arg
		}
		results := SearchZones(strings.ToLower(q))
		if len(results) < 1 {
			fmt.Fprintf(os.Stderr(), "Unknown time zone %s\n", q)
			os.Exit(3)
			return nil
		}
		results.Print(os.Stdout())
		os.Exit(0)
		return nil
	}

	config, err := LoadDefaultConfig(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr(), "Config error: %s\n", err)
		os.Exit(2)
		return nil
	}

	var initialModel = model{
		zones:      config.Zones,
		keymaps:    config.Keymaps,
		clock:      *NewClockNow(),
		showDates:  false,
		isMilitary: *military,
		watch:      *watch,
		showHelp:   false,
		zoneStyle:  AbbreviationZoneStyle,
	}

	if when != nil {
		initialModel.clock = *NewClockUnixTimestamp(*when)
	}

	initialModel.interactive = !*exitQuick && isatty.IsTerminal(os.Stdout().Fd())

	return &initialModel
}

func main() {
	SetupLogger()
	initialModel := parseMainArgs()
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr(), "Alas, there's been an error: %v", err)
		os.Exit(1)
		return
	}
}
