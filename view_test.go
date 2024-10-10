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
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/tools/txtar"
)

var (
	dstTestZones []*Zone = nil
)

func LoadDstTestZones(t *testing.T) []*Zone {
	if dstTestZones == nil {
		config, err := LoadDefaultConfig([]string{
			"UTC",                              // Z
			"Europe/Paris",                     // Z+1 (CET), Z+2 (CEST)
			"Israel",                           // Z+2 (IST), Z+3 (IDT)
			"Asia/Calcutta",                    // Z+5:30 (IST)
			"Australia/Eucla",                  // Z+8:45 (no abbreviation)
			"Australia/Sydney",                 // Z+10 (AEST), Z+11 (AEDT)
			"Pacific/Kiritimati",               // Z+14 (no abbreviation)
			"Pacific/Honolulu",                 // Z-10 (HST)
			"US/Central",                       // Z-6 (CST), Z-5 (CDT)
			"Cuba",                             // Z-5 (CST), Z-4 (CDT)
			"America/Argentina/ComodRivadavia", // Z-3 (no abbreviation)
		})
		if err != nil {
			t.Fatal(err)
		}
		dstTestZones = config.Zones[1:]
	}
	return dstTestZones
}

func RunDstDaysTest(
	t *testing.T,
	title string,
	testDataFile string,
	transition time.Time,
) {
	testData, err := txtar.ParseFile(testDataFile)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		time time.Time
	}{
		{title, transition},
		{"Hour before", transition.Add(-time.Hour)},
		{"Hour after", transition.Add(time.Hour)},
		{"3 days before", transition.Add(-72 * time.Hour)},
		{"2 days before", transition.Add(-48 * time.Hour)},
		{"1 day before", transition.Add(-24 * time.Hour)},
		{"Day after", transition.Add(24 * time.Hour)},
	}

	var observations []string
	var outputData = []txtar.File{}
	for _, test := range tests {
		state := model{
			zones:      LoadDstTestZones(t),
			clock:      *NewClockTime(test.time),
			isMilitary: true,
			showDates:  true,
		}
		observed := stripAnsiControlSequences(state.View())
		observations = append(observations, observed)
		outputData = append(
			outputData,
			txtar.File{
				Name: fmt.Sprintf("%v (%v = %v)", test.name, test.time.Format(time.RFC3339), test.time.Unix()),
				Data: []byte(observed),
			},
		)
	}

	archive := txtar.Archive{
		Comment: testData.Comment,
		Files: outputData,
	}
	os.WriteFile(testDataFile, txtar.Format(&archive), 0666)

	for i, test := range tests {
		expected := stripAnsiControlSequencesAndNewline(testData.Files[i].Data)
		observed := stripAnsiControlSequencesAndNewline(outputData[i].Data)
		if expected != observed {
			t.Errorf("Daylight Saving: Mismatched %s: Check git diff %s", test.name, testDataFile)
		}
	}
}

func TestDstEndDays(t *testing.T) {
	testDataFile := "testdata/view/test-dst-end-days.txt"
	europeEndDst := time.Date(2024, time.October, 27, 1, 0, 0, 0, time.UTC)
	RunDstDaysTest(t, "Europe DST end", testDataFile, europeEndDst)
}

func TestDstStartDays(t *testing.T) {
	testDataFile := "testdata/view/test-dst-start-days.txt"
	europeStartDst := time.Date(2024, time.March, 31, 1, 0, 0, 0, time.UTC)
	RunDstDaysTest(t, "Europe DST start", testDataFile, europeStartDst)
}

func TestDstSpecialMidnights(t *testing.T) {
	// The following are expressed in UTC, because the local date boundary is unusual:
	cubaDstStart := time.Date(2017, time.March, 12, 5, 0, 0, 0, time.UTC) // 12 Mar Cuba missing midnight
	cubaDstEnd := time.Date(2017, time.November, 5, 4, 0, 0, 0, time.UTC) // 5 Mar Cuba double midnight

	testDataFile := "testdata/view/test-dst-midnights.txt"
	testData, err := txtar.ParseFile(testDataFile)
	if err != nil {
		t.Fatal(err)
	}

	config, err := LoadDefaultConfig([]string{
		"UTC",  // Z
		"Cuba", // Z-5 (CST), Z-4 (CDT)
	})
	if err != nil {
		t.Fatal(err)
	}
	midnightTestZones := config.Zones[1:]

	tests := []struct {
		name string
		time time.Time
	}{
		{"Start DST missing midnight", cubaDstStart},
		{"End DST double midnight", cubaDstEnd},
	}

	var outputData = make([]txtar.File, len(tests))
	for i, test := range tests {
		state := model{
			zones:      midnightTestZones,
			clock:      *NewClockTime(test.time),
			isMilitary: true,
			showDates:  true,
		}
		observed := stripAnsiControlSequences(state.View())
		outputData[i] = txtar.File{
			Name: fmt.Sprintf("%v (%v = %v)", test.name, test.time.Format(time.RFC3339), test.time.Unix()),
			Data: []byte(observed),
		}
	}

	archive := txtar.Archive{
		Comment: testData.Comment,
		Files: outputData,
	}
	os.WriteFile(testDataFile, txtar.Format(&archive), 0666)

	for i, test := range tests {
		observed := stripAnsiControlSequencesAndNewline(outputData[i].Data)
		expected := stripAnsiControlSequencesAndNewline(testData.Files[i].Data)
		if observed != expected {
			t.Errorf("Midnight DST: Mismatched %s: Check git diff %s", test.name, testDataFile)
		}
	}
}

func TestFractionalTimezoneOffsets(t *testing.T) {
	testDataFile := "testdata/view/test-fractional-timezone-offsets.txt"
	testData, err := txtar.ParseFile(testDataFile)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		datetime string
		keystrokes string
	}{
		{"Start time", "2017-11-05T00:29:02Z", ""},
		{"Plus 1 hour", "2017-11-05T01:29:02Z", "l"},
		{"Minus 15 minutes", "2017-11-05T01:14:02Z", "---------------"}, // @todo implement new keys
		{"Return to start", "2017-11-05T00:29:02Z", "HL<>+++++++++++++++h"},
		{"Minus 1 hour, date changed", "2017-11-04T23:29:02Z", "h"},
		{"Return to start", "2017-11-05T00:29:02Z", "l"},
	}

	start, err := time.Parse(time.RFC3339, tests[0].datetime)
	if err != nil {
		t.Fatalf("Could not parse test Start time: %v", err)
	}

	state := &model{
		zones:      LoadDstTestZones(t),
		clock:      *NewClockTime(start),
		keymaps:    DefaultKeymaps,
		isMilitary: true,
		showDates:  true,
	}

	var observations []string
	var outputData = []txtar.File{}
	for _, test := range tests {
		for k, key := range test.keystrokes {
			msg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{key},
				Alt:   false,
			}

			_, cmd := state.Update(msg)
			if cmd != nil {
				t.Fatalf("Expected nil Cmd for '%v' (key %v), but got %v", key, k, cmd)
			}
		}

		observed := stripAnsiControlSequences(state.View())
		observations = append(observations, observed)
		observedDatetime := state.clock.t.Format(time.RFC3339)
		observedUnixtime := state.clock.t.Unix()
		if observedDatetime != test.datetime {
			t.Errorf("Fraction Timezones: Mismatched datetime for %v: expected %v but got %v", test.name, test.datetime, observedDatetime)
		}
		outputData = append(
			outputData,
			txtar.File{
				Name: fmt.Sprintf("%v (%v = %v)", test.name, observedDatetime, observedUnixtime),
				Data: []byte(observed),
			},
		)
	}

	archive := txtar.Archive{
		Comment: testData.Comment,
		Files: outputData,
	}
	os.WriteFile(testDataFile, txtar.Format(&archive), 0666)

	for i, test := range tests {
		var expected string = ""
		if len(testData.Files) > 0 {
			expected = stripAnsiControlSequencesAndNewline(testData.Files[i].Data)
		}
		observed := stripAnsiControlSequencesAndNewline(outputData[i].Data)
		if expected != observed {
			t.Errorf("Fraction Timezones: Mismatched %s: Check git diff %s", test.name, testDataFile)
		}
	}
}

func TestRightAlignment(t *testing.T) {
	testDataFile := "testdata/view/test-right-alignment.txt"
	testData, err := txtar.ParseFile(testDataFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := stripAnsiControlSequencesAndNewline(testData.Files[0].Data)

	tests := []struct {
		columns int
	}{
		{columns: 1},
		{columns: 80},
		{columns: 999},
	}

	originalColumns := os.Getenv("COLUMNS")
	var observations []string
	for _, test := range tests {
		os.Setenv("COLUMNS", fmt.Sprintf("%v", test.columns))
		observed := stripAnsiControlSequences(utcMinuteAfterMidnightModel.View())
		observations = append(observations, observed)
	}
	os.Setenv("COLUMNS", originalColumns)

	archive := txtar.Archive{
		Comment: testData.Comment,
		Files: []txtar.File{
			{
				Name: "expected",
				Data: []byte(expected),
			},
			{
				Name: "observed: narrow",
				Data: []byte(observations[0]),
			},
			{
				Name: "observed: medium",
				Data: []byte(observations[1]),
			},
			{
				Name: "observed: wide",
				Data: []byte(observations[2]),
			},
		},
	}
	os.WriteFile(testDataFile, txtar.Format(&archive), 0666)

	expectations := strings.Split(expected, "\n")
	for i, test := range tests {
		if !strings.Contains(observations[i], expectations[i]) {
			t.Errorf("Expected %d-column alignment “%s”, but got: “%s”", test.columns, expectations[i], observations[i])
		}
	}
}
