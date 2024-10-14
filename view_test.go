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
		config, err := LoadConfig([]string{
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

	config, err := LoadConfig([]string{
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
		{"Minus 15 minutes", "2017-11-05T01:14:02Z", "---------------"},
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

func TestHighlightMarkers(t *testing.T) {
	testDataFile := "testdata/view/test-highlight-markers.txt"
	testData, err := txtar.ParseFile(testDataFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := stripAnsiControlSequencesAndNewline(testData.Files[0].Data)

	keys := "jkj" // down, up, down

	var state = utcMinuteAfterMidnightModel
	for k, key := range keys {
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

	observed := txtar.File{
		Name: "Highlight Local Zone",
		Data: []byte(stripAnsiControlSequences(state.View())),
	}

	archive := txtar.Archive{
		Comment: testData.Comment,
		Files: []txtar.File{observed},
	}
	os.WriteFile(testDataFile, txtar.Format(&archive), 0666)

	if expected != stripAnsiControlSequencesAndNewline(observed.Data) {
		t.Errorf("Fraction Timezones: Mismatched highlight markers: Check git diff %s", testDataFile)
	}
}

// Test all vertical alignments, from the perspectives of different local zones
func TestLocalTimezones(t *testing.T) {
	testDataFile := "testdata/view/test-local-timezones.txt"
	testData, err := txtar.ParseFile(testDataFile)
	if err != nil {
		t.Fatal(err)
	}

	testGroups := [][]struct {
		localZone string
		localTime string
	}{
		{
			{"UTC", "2017-11-05T00:00:00Z"},
			{"Asia/Calcutta", "2017-11-05T05:30:00+05:30"},
			{"Cuba", "2017-11-04T20:00:00-04:00"},
		},
		{
			{"UTC", "2017-11-06T00:30:00Z"},
			{"Asia/Calcutta", "2017-11-06T06:00:00+05:30"},
			{"Cuba", "2017-11-05T19:30:00-05:00"},
		},
	}

	displayZones := LoadDstTestZones(t)
	var outputData = []txtar.File{}
	for i, testGroup := range testGroups {
		for j, test := range testGroup {
			testTime, err := time.Parse(time.RFC3339, test.localTime)
			if err != nil {
				t.Fatalf("Could not parse test time configuration [%v][%v]: '%v'", i, j, test.localTime)
			}

			localZoneAtTop := make([]*Zone, 0, len(displayZones) + 1)
			for _, zone := range(displayZones) {
				if zone.DbName == test.localZone {
					// Copy localZone to the top of list to render all other zones relative to it
					localZone := *zone
					localZone.Name = "Local"
					localZoneAtTop = append(localZoneAtTop, &localZone)
					localZoneAtTop = append(localZoneAtTop, displayZones...)
					break
				}
			}
			if len(localZoneAtTop) == 0 {
				t.Fatalf("Could not find displayable timezone for case [%v][%v]: '%v'", i, j, test.localTime)
			}

			state := model{
				zones:       localZoneAtTop,
				clock:       *NewClockTime(testTime),
				isMilitary:  true,
				showDates:   true,
				formatStyle: IsoFormatStyle,
				zoneStyle:   WithRelativeZoneStyle,
			}

			observed := stripAnsiControlSequences(state.View())
			outputData = append(outputData, txtar.File{
				Name: fmt.Sprintf("[%v][%v] %v (%v = %v)", i, j, test.localZone, testTime.Format(time.RFC3339), testTime.Unix()),
				Data: []byte(observed),
			})
		}
	}

	archive := txtar.Archive{
		Comment: testData.Comment,
		Files: outputData,
	}
	os.WriteFile(testDataFile, txtar.Format(&archive), 0666)

	var count = 0
	for i, testGroup := range testGroups {
		// Implementation explained in the testData file:
		comparisonColumns := make([]string, len(testGroup))

		for j, test := range testGroup {
			// Test for any changes
			observed := stripAnsiControlSequencesAndNewline(outputData[count].Data)
			expected := stripAnsiControlSequencesAndNewline(testData.Files[count].Data)
			if observed != expected {
				t.Errorf("Local Timezones: Unexpected result [%v][%v] for %v: Check git diff %v", i, j, test.localZone, testDataFile)
			}

			// Test for expected properties (this is sensitive to the layout format)
			hourIndex := 11
			localTime := test.localTime[hourIndex:hourIndex + 2]
			if localTime[0] == '0' {
				localTime = " " + localTime[1:]
			}
			localTime = " " + localTime + " "

			localLine := 4
			lines := strings.Split(observed, "\n")
			columnIndex := strings.Index(lines[localLine], localTime)
			if columnIndex < 0 {
				t.Errorf("Local Timezones: Failed [%v][%v] for %v: Could not find local hour %v", i, j, test.localZone, localTime)
			} else {
				rowGap := 3
				column := strings.Builder{}
				for l := localLine + rowGap; l < len(lines); l += rowGap {
					line := lines[l]
					column.WriteString(line[columnIndex:columnIndex + 3])
				}
				comparisonColumns[j] = column.String()
			}

			count = count + 1
		}

		for _, result := range comparisonColumns {
			if result != comparisonColumns[0] {
				t.Errorf("Local Timezones: Inconsistencies in group [%v]:\n%v", i, strings.Join(comparisonColumns, "\n"))
				break
			}
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

func TestTimeFormats(t *testing.T) {
	testDataFile := "testdata/view/test-time-formats.txt"
	testData, err := txtar.ParseFile(testDataFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := stripAnsiControlSequencesAndNewline(testData.Files[0].Data)

	tests := []struct {
		name string
		formatStyle FormatStyle
	}{
		{"DefaultFormatStyle", DefaultFormatStyle},
		{"IsoFormatStyle", IsoFormatStyle},
		{"UnixFormatStyle", UnixFormatStyle},
	}

	var observations []string
	var outputs = []txtar.File{
		{
			Name: "expected",
			Data: []byte(expected),
		},
	}
	oldHasDarkBackground := hasDarkBackground
	hasDarkBackground = true
	var state = utcMinuteAfterMidnightModel
	for i, test := range tests {
		if state.formatStyle != test.formatStyle {
			t.Errorf("Expected %s %v for test %d but got: %v", test.name, test.formatStyle, i, state.formatStyle)
		}
		observed := stripAnsiControlSequences(state.View())
		observations = append(observations, observed)
		outputs = append(
			outputs,
			txtar.File{
				Name: fmt.Sprintf("observed: %v", test.name),
				Data: []byte(observed),
			},
		)
		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'f'},
			Alt:   false,
		}
		_, cmd := state.Update(msg)
		if cmd != nil {
			t.Fatalf("Expected nil Cmd, but got %v", cmd)
		}
	}
	hasDarkBackground = oldHasDarkBackground

	for i := len(tests) - 1; i >= 0; i-- {
		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'F'},
			Alt:   false,
		}
		_, cmd := state.Update(msg)
		if cmd != nil {
			t.Fatalf("Expected nil Cmd, but got %v", cmd)
		}
		if state.formatStyle != tests[i].formatStyle {
			t.Errorf("Expected %s %v for reverse test %d but got: %v", tests[i].name, tests[i].formatStyle, i, state.formatStyle)
		}
	}

	archive := txtar.Archive{
		Comment: testData.Comment,
		Files: outputs,
	}
	os.WriteFile(testDataFile, txtar.Format(&archive), 0666)

	expectations := strings.Split(expected, "\n")
	for i, test := range tests {
		if !strings.Contains(observations[i], expectations[i]) {
			t.Errorf("Expected %v “%s”, but got: “%s”", test.name, expectations[i], observations[i])
		}
	}
}

func TestZoneStyles(t *testing.T) {
	testDataFile := "testdata/view/test-zone-styles.txt"
	testData, err := txtar.ParseFile(testDataFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := stripAnsiControlSequencesAndNewline(testData.Files[0].Data)

	tests := []struct {
		name string
		zoneStyle ZoneStyle
	}{
		{"AbbreviationZoneStyle", AbbreviationZoneStyle},
		{"WithZOffsetZoneStyle", WithZOffsetZoneStyle},
		{"WithRelativeZoneStyle", WithRelativeZoneStyle},
	}

	var observations []string
	var outputs = []txtar.File{
		{
			Name: "expected",
			Data: []byte(expected),
		},
	}
	oldHasDarkBackground := hasDarkBackground
	hasDarkBackground = true
	var state = utcMinuteAfterMidnightModel
	state.isMilitary = false
	for i, test := range tests {
		if state.zoneStyle != test.zoneStyle {
			t.Errorf("Expected %s %v for test %d but got: %v", test.name, test.zoneStyle, i, state.zoneStyle)
		}
		observed := stripAnsiControlSequences(state.View())
		observations = append(observations, observed)
		outputs = append(
			outputs,
			txtar.File{
				Name: fmt.Sprintf("observed: %v", test.name),
				Data: []byte(observed),
			},
		)
		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'z'},
			Alt:   false,
		}
		_, cmd := state.Update(msg)
		if cmd != nil {
			t.Fatalf("Expected nil Cmd, but got %v", cmd)
		}
	}
	hasDarkBackground = oldHasDarkBackground

	for i := len(tests) - 1; i >= 0; i-- {
		msg := tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'Z'},
			Alt:   false,
		}
		_, cmd := state.Update(msg)
		if cmd != nil {
			t.Fatalf("Expected nil Cmd, but got %v", cmd)
		}
		if state.zoneStyle != tests[i].zoneStyle {
			t.Errorf("Expected %s %v for reverse test %d but got: %v", tests[i].name, tests[i].zoneStyle, i, state.zoneStyle)
		}
	}

	archive := txtar.Archive{
		Comment: testData.Comment,
		Files: outputs,
	}
	os.WriteFile(testDataFile, txtar.Format(&archive), 0666)

	expectations := strings.Split(expected, "\n")
	for i, test := range tests {
		if !strings.Contains(observations[i], expectations[i]) {
			t.Errorf("Expected %v “%s”, but got: “%s”", test.name, expectations[i], observations[i])
		}
	}
}
