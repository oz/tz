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
	"os"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
)

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
		os.Setenv("COLUMNS", strconv.Itoa(test.columns))
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
