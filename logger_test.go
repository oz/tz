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
	"math/rand"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	oldDebug := os.Getenv("DEBUG")
	os.Setenv("DEBUG", "1")
	SetupLogger()

	logMsg := fmt.Sprintf("Hello world %v!", rand.Int())
	logger.Printf(logMsg)

	out, err := os.ReadFile("debug.log")
	if err != nil {
		t.Errorf("Could not read log file debug.log: %v", err)
	}
	lines := strings.Split(string(out), "\n")
	lastLine := lines[len(lines) - 2] // ignore trailing newline
	if !strings.Contains(lastLine, logMsg) {
		t.Errorf("Missing log line in debug.log: %s", logMsg)
	}

	if oldDebug != "1" {
		os.Setenv("DEBUG", oldDebug)
		SetupLogger()
	}
}

func TestNoLogger(t *testing.T) {
	oldDebug, debugWasSet := os.LookupEnv("DEBUG")
	os.Unsetenv("DEBUG")
	SetupLogger()

	logMsg := fmt.Sprintf("Not logged %v!", rand.Int())
	logger.Print(logMsg)

	out, err := os.ReadFile("debug.log")
	if err != nil {
		t.Errorf("Could not read log file debug.log: %v", err)
	}
	if strings.Contains(string(out), logMsg) {
		t.Errorf("Log file debug.log contained forbidden string: %s", logMsg)
	}

	if debugWasSet {
		os.Setenv("DEBUG", oldDebug)
		SetupLogger()
	}
}
