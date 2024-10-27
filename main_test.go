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
	"regexp"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"golang.org/x/tools/txtar"
)

var (
	ansiControlSequenceRegexp = regexp.MustCompile(regexp.QuoteMeta(termenv.CSI) + "[^m]*m")
	utcMinuteAfterMidnightTime = time.Date(
		2017, // Year
		11, // Month
		5, // Day
		0, // Hour
		1, // Minutes
		2, // Seconds
		127, // Nanoseconds
		time.UTC,
	)
	utcMinuteAfterMidnightModel = model{
		zones:      DefaultZones[len(DefaultZones) - 1:],
		clock:      *NewClockTime(utcMinuteAfterMidnightTime),
		keymaps:    DefaultKeymaps,
		isMilitary: true,
		showDates:  true,
	}
)

func failUnlessExpectedError(t *testing.T, err error, expectedError string, contextFormat string, a ...any) {
	msg := fmt.Sprintf(contextFormat, a...)
	if err != nil {
		if !strings.Contains(err.Error(), expectedError) {
			t.Fatalf("Expected specific error %s, but got a different error: %v", msg, err)
		}
	} else {
		t.Fatalf("Expected error %s, but none occurred", msg)
	}
}

func getTimestampWithHour(hour int) time.Time {
	return time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		hour,
		43, // Minutes
		59, // Seconds
		127, // Nanoseconds
		time.Now().Location(),
	)
}

func parseMainArgsWithPanicRecovery(panicErr *error) *model {
	*panicErr = nil
	savedFlags := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet("main_test", flag.PanicOnError)

	defer func() {
		flag.CommandLine = savedFlags
		if err, recovered := recover().(error); recovered {
			*panicErr = err
		}
	}()

	return parseMainArgs()
}

func stripAnsiControlSequences(s string) string {
	return ansiControlSequenceRegexp.ReplaceAllString(s, "")
}

func stripAnsiControlSequencesAndNewline(bytes []byte) string {
	s := strings.TrimSuffix(string(bytes), "\n")
	return ansiControlSequenceRegexp.ReplaceAllString(s, "")
}

func testMainArgWhen(t *testing.T, when string, whenSeconds int64) {
	var err error
	var osWrapper = NewTestingOsWrapper(t)

	// 1. Test initial state with -when (overrides -w)
	osWrapper.HomeDir = "."
	osWrapper.Setargs([]string{"-when", when, "-w"})
	model := parseMainArgsWithPanicRecovery(&err)

	if err != nil {
		t.Errorf("Unexpected failure for `-when %v`: %v", when, err)
	}
	if osWrapper.ExitCode != nil {
		t.Errorf("Unexpected exit code for -when %v`: %v", when, *osWrapper.ExitCode)
	}
	if model == nil {
		t.Fatalf("Model was nil after parsing `-when %v`", when)
	}
	if model.clock.t.Unix() != whenSeconds {
		t.Errorf("Model `-when %v` clock.time was incorrect: %v (%v)", when, model.clock.t.Unix(), model.clock.t)
	}
	if model.clock.isRealTime {
		t.Errorf("Model `-when %v` clock.isRealTime was incorrect: %v", when, model.clock.isRealTime)
	}

	// 2. Test that tick does not change the state with -when
	tickMsg := tickMsg(time.Now())
	if _, cmd := model.Update(tickMsg); cmd == nil {
		t.Fatalf("Expected non-nil Cmd, but got %v", cmd)
	}
	if model.clock.t.Unix() != whenSeconds {
		t.Errorf("Model `-when %v` clock.time was unstable after tickMsg: %v (%v)", when, model.clock.t.Unix(), model.clock.t)
	}
	if model.clock.isRealTime {
		t.Errorf("Model `-when %v` clock.isRealTime was incorrect after tickMsg: %v", when, model.clock.isRealTime)
	}

	// 3. Cancel `-when` and activate `-w` using the interactive `t` key
	keyMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'t'},
	}
	if _, cmd := model.Update(keyMsg); cmd != nil {
		t.Fatalf("Expected nil Cmd, but got %v", cmd)
	}
	if model.clock.t.Unix() == whenSeconds {
		t.Errorf("Model clock.time was incorrect after `t` key: %v (time %v)", model.clock.t.Unix(), model.clock.t)
	}
	if !model.clock.isRealTime {
		t.Errorf("Model clock.isRealTime was incorrect after `t` key: %v", model.clock.isRealTime)
	}

	// 4. Test that ticks are tracking the current time with -w
	oldTime := model.clock.t
	if _, cmd := model.Update(tickMsg); cmd == nil {
		t.Fatalf("Expected non-nil Cmd, but got %v", cmd)
	}
	if model.clock.t == oldTime {
		t.Errorf("Model clock.time was incorrect with -w option: %v", model.clock.t)
	}
	if !model.clock.isRealTime {
		t.Errorf("Model clock.isRealTime was incorrect with -w option: %v", model.clock.isRealTime)
	}
}

func TestMain(m *testing.M) {
	SetupLogger()
	m.Run()
}

func TestMainArgNone(t *testing.T) {
	var err error
	var osWrapper = NewTestingOsWrapper(t)

	osWrapper.HomeDir = "."
	osWrapper.Setargs([]string{})
	model := parseMainArgsWithPanicRecovery(&err)

	if err != nil {
		t.Errorf("Unexpected failure for empty flag args: %v", err)
	}
	if osWrapper.ExitCode != nil {
		t.Errorf("Unexpected exit code for zero flag args: %v", *osWrapper.ExitCode)
	}
	if model == nil {
		t.Errorf("Model was nil after parsing flag args")
	} else {
		if model.isMilitary {
			t.Errorf("Default model.isMilitary was %v but expected false", model.isMilitary)
		}
		if model.showDates {
			t.Errorf("Default model.showDates was %v but expected false", model.showDates)
		}
		if model.showHelp {
			t.Errorf("Default model.showHelp was %v but expected false", model.showHelp)
		}
		if model.watch {
			t.Errorf("Default model.watch was %v but expected false", model.watch)
		}
		if !model.clock.isRealTime {
			t.Errorf("Default model.clock.isRealTime was %v but expected true", model.clock.isRealTime)
		}
	}

	osWrapper.HomeDir = ""
	expectedOutput := "Config error: File error: TestingOsWrapper.HomeDir is not yet set"
	parseMainArgsWithPanicRecovery(&err)
	if osWrapper.ExitCode == nil {
		t.Error("Main should exit for invalid HomeDir, but it did not")
	} else if *osWrapper.ExitCode != 2 {
		t.Errorf("Main should exit with code 0 for invalid HomeDir, but got %v", *osWrapper.ExitCode)
	}
	if output := strings.TrimSpace(osWrapper.ConsumeStderr()); output != expectedOutput {
		t.Errorf("Main should have printed '%v', but got '%v'", expectedOutput, output)
	}
}

func TestMainArgInvalid(t *testing.T) {
	var err error
	var osWrapper = NewTestingOsWrapper(t)

	osWrapper.Setargs([]string{"- "})
	osWrapper.RedirectPlatformOutput() // capture "flag" package messages
	parseMainArgsWithPanicRecovery(&err)
	osWrapper.RevertPlatformOutput() // print "testing" package messages
	failUnlessExpectedError(t, err, "flag provided but not defined: - ", "in %s", t.Name())

	expectedOutputFragment := "flag provided but not defined: - \nUsage"
	if output := strings.TrimSpace(osWrapper.ConsumeStderr()); !strings.Contains(output, expectedOutputFragment) {
		t.Errorf("Main `- ` flag arg should have printed '%v', but got '%v'", expectedOutputFragment, output)
	}
}

func TestMainArgList(t *testing.T) {
	var err error
	var osWrapper = NewTestingOsWrapper(t)

	osWrapper.Setargs([]string{"-list", "UTC"})
	parseMainArgsWithPanicRecovery(&err)
	if err != nil {
		t.Errorf("Unexpected failure for -list flag: %v", err)
	}
	if osWrapper.ExitCode == nil {
		t.Error("Main -list should exit, but it did not")
	} else if *osWrapper.ExitCode != 0 {
		t.Errorf("Main -list should exit with code 0, but got %v", *osWrapper.ExitCode)
	}

	stderr := osWrapper.ConsumeStderr()
	if len(stderr) > 0 {
		t.Errorf("Unexpected stderr for -list flag: %v", stderr)
	}

	stdout := osWrapper.ConsumeStdout()
	expectedOutputFragment := "UTC (+00:00) :: UTC"
	if !strings.Contains(stdout, expectedOutputFragment) {
		t.Errorf("Stdout for -list flag did not contain '%v': %v", expectedOutputFragment, stdout)
	}

	osWrapper.Setargs([]string{"-list", "!"})
	expectedErrOutput := "Unknown time zone !"
	parseMainArgsWithPanicRecovery(&err)
	if err != nil {
		t.Errorf("Unexpected failure for `-list !` flag: %v", err)
	}
	if osWrapper.ExitCode == nil {
		t.Error("Main `-list !` should exit, but did not")
	} else if *osWrapper.ExitCode != 3 {
		t.Errorf("Main `-list !` should exit with code 3, but got %v", *osWrapper.ExitCode)
	}
	if errOutput := strings.TrimSpace(osWrapper.ConsumeStderr()); errOutput != expectedErrOutput {
		t.Errorf("Main `-list !` should have stderr '%v', but got '%v'", expectedErrOutput, errOutput)
	}
	if stdOutput := osWrapper.ConsumeStdout(); len(stdOutput) > 0 {
		t.Errorf("Main `-list !` should have empty stdout, but got '%v'", stdOutput)
	}
}

func TestMainArgVersion(t *testing.T) {
	var err error
	var osWrapper = NewTestingOsWrapper(t)

	osWrapper.Setargs([]string{"-v"})
	parseMainArgsWithPanicRecovery(&err)

	if err != nil {
		t.Errorf("Unexpected failure for -v flag: %v", err)
	}

	stderr := osWrapper.ConsumeStderr()
	if len(stderr) > 0 {
		t.Errorf("Unexpected stderr for -v flag: %v", stderr)
	}

	stdout := strings.TrimSpace(osWrapper.ConsumeStdout())
	expectedOutput := fmt.Sprintf("tz %v", CurrentVersion)
	if stdout != expectedOutput {
		t.Errorf("Unexpected stdout for -v flag (expected '%v'): %v", expectedOutput, stdout)
	}

	if osWrapper.ExitCode == nil {
		t.Error("Main -v should exit, but did not")
	} else if *osWrapper.ExitCode != 0 {
		t.Errorf("Main -v should exit with code 0, but got %v", *osWrapper.ExitCode)
	}
}

func TestMainArgWhen(t *testing.T) {
	testMainArgWhen(t, "1", 1)
	testMainArgWhen(t, "-1", -1)
	testMainArgWhen(t, "0", 0)
	testMainArgWhen(t, "1970-01-01T00:00:01Z", 1)
	testMainArgWhen(t, "2006-01-02T15:04:05-07:00", 1136239445)
}

func TestMainArgWhenInvalid(t *testing.T) {
	var err error
	var osWrapper = NewTestingOsWrapper(t)

	osWrapper.Setargs([]string{"-when", "midnight"})
	osWrapper.RedirectPlatformOutput() // capture "flag" package messages
	parseMainArgsWithPanicRecovery(&err)
	osWrapper.RevertPlatformOutput() // print "testing" package messages
	failUnlessExpectedError(t, err, "invalid value \"midnight\" for flag -when: Could not parse", "in %s", t.Name())

	expectedUsageFragment := "date-time in seconds since unix epoch, or in ISO8601/RFC3339 format"
	if output := strings.TrimSpace(osWrapper.ConsumeStderr()); !strings.Contains(output, expectedUsageFragment) {
		t.Errorf("Main `-when midnight` flag arg should have printed '%v', but got '%v'", expectedUsageFragment, output)
	}
}

func TestUpdateIncHour(t *testing.T) {
	// "l" key -> go right
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'l'},
		Alt:   false,
	}

	tests := []struct {
		startHour          int
		nextHour           int
		changesClockDateBy int
	}{
		{startHour: 0, nextHour: 1, changesClockDateBy: 0},
		{startHour: 1, nextHour: 2, changesClockDateBy: 0},
		// ...
		{startHour: 23, nextHour: 0, changesClockDateBy: 1},
	}

	for _, test := range tests {
		m := model{
			zones:   DefaultZones,
			keymaps: DefaultKeymaps,
			clock:   *NewClockTime(getTimestampWithHour(test.startHour)),
		}

		tb := m.clock.Time()
		db := tb.Day()
		_, cmd := m.Update(msg)
		ta := m.clock.Time()
		da := ta.Day()

		if cmd != nil {
			t.Fatalf("Expected nil Cmd, but got %v", cmd)
		}
		h := m.clock.t.Hour()
		if h != test.nextHour {
			t.Errorf("Expected %d, but got %d", test.nextHour, h)
		}
		if test.changesClockDateBy != 0 && da == db {
			t.Errorf("Expected date change of %d day, but got %d", test.changesClockDateBy, da-db)
		}
		if ta.Minute() != tb.Minute() {
			t.Errorf("Unexpected change of minute from '%s' to '%s'", tb, ta)
		}
	}
}

func TestUpdateDecHour(t *testing.T) {
	// "h" key -> go left
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'h'},
		Alt:   false,
	}

	tests := []struct {
		startHour int
		nextHour  int
	}{
		{startHour: 23, nextHour: 22},
		{startHour: 22, nextHour: 21},
		// ...
		{startHour: 0, nextHour: 23},
	}

	for _, test := range tests {
		m := model{
			zones:   DefaultZones,
			keymaps: DefaultKeymaps,
			clock:   *NewClockTime(getTimestampWithHour(test.startHour)),
		}
		_, cmd := m.Update(msg)
		if cmd != nil {
			t.Fatalf("Expected nil Cmd, but got %v", cmd)
		}
		h := m.clock.t.Hour()
		if h != test.nextHour {
			t.Errorf("Expected %d, but got %d", test.nextHour, h)
		}
	}
}

func TestUpdateShowDatesMsg(t *testing.T) {
	// "d" key -> toggle dates
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'d'},
	}

	m := utcMinuteAfterMidnightModel

	dateMarker := "📆"

	if !strings.Contains(m.View(), dateMarker) {
		t.Fatalf("Dates should be shown in utcMinuteAfterMidnightModel")
	}

	if _, cmd := m.Update(msg); cmd != nil {
		t.Fatalf("Expected nil Cmd, but got %v", cmd)
	}

	if strings.Contains(m.View(), dateMarker) {
		t.Fatalf("Dates should have been toggled after `d` key")
	}
}

func TestUpdateShowHelpMsg(t *testing.T) {
	// "?" key -> help
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'?'},
		Alt:   false,
	}

	m := utcMinuteAfterMidnightModel
	if m.showHelp {
		t.Error("showHelp should be disabled by default")
	}

	status1 := status(m)

	if _, cmd := m.Update(msg); cmd != nil {
		t.Fatalf("Expected nil Cmd, but got %v", cmd)
	}

	if !m.showHelp {
		t.Error("showHelp not enabled by '?' key")
	}

	oldHasDarkBackground := hasDarkBackground
	hasDarkBackground = true
	status2 := status(m)
	hasDarkBackground = oldHasDarkBackground

	if status2 == status1 || !strings.Contains(status2, "d: toggle date") {
		t.Errorf("Expected help, but got:\n%v", status2)
	}

	if _, cmd := m.Update(msg); cmd != nil {
		t.Fatalf("Expected nil Cmd, but got %v", cmd)
	}

	if m.showHelp {
		t.Error("showHelp not toggled")
	}

	status3 := status(m)

	if status3 != status1 {
		t.Errorf("Expected final status identical to initial status, but got:\n%v", status3)
	}
}

func TestUpdateQuitMsg(t *testing.T) {
	// "q" key -> quit
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
		Alt:   false,
	}

	m := model{
		zones:   DefaultZones,
		keymaps: DefaultKeymaps,
		clock:   *NewClockTime(getTimestampWithHour(10)),
	}
	_, cmd := m.Update(msg)
	if cmd == nil {
		t.Fatalf("Expected tea.Quit Cmd, but got %v", cmd)
	}
	// tea.Quit is a function, we can't really test with == here, and
	// calling it is getting into internal territory.
}

func TestMilitaryTime(t *testing.T) {
	testDataFile := "testdata/main/test-military-time.txt"
	testData, err := txtar.ParseFile(testDataFile)
	if err != nil {
		t.Fatal(err)
	}

	formatted := utcMinuteAfterMidnightTime.Format(" 15:04, Mon Jan 02, 2006")
	expected := stripAnsiControlSequencesAndNewline(testData.Files[0].Data)
	observed := stripAnsiControlSequences(utcMinuteAfterMidnightModel.View())

	archive := txtar.Archive{
		Comment: testData.Comment,
		Files: []txtar.File{
			{
				Name: "expected",
				Data: []byte(expected),
			},
			{
				Name: "observed",
				Data: []byte(observed),
			},
		},
	}
	os.WriteFile(testDataFile, txtar.Format(&archive), 0666)

	if formatted != expected {
		t.Errorf("Expected military time of %s, but got %s", expected, formatted)
	}
	if !strings.Contains(observed, expected) {
		t.Errorf("Expected military time of %s, but got %s", expected, observed)
	}
}
