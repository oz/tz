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
	"strings"
	"testing"
	"time"
)

func TestLoadConfigParser(t *testing.T) {
	tests := []struct {
		args     []string
		envArg   string
		expected string
	}{
		{nil, "", "Local;UTC"},
		{[]string{"GMT"}, "UTC", "Local;GMT"},
		{[]string{"", "GMT", ""}, "UTC", "Local;UTC;GMT;UTC"},
		{nil, "GMT", "Local;GMT"},
		{nil, ";GMT;", "Local;UTC;GMT;UTC"},
		{nil, "Unknown", ""},
	}

	oldEnv := os.Getenv("TZ_LIST")
	for _, test := range tests {
		os.Setenv("TZ_LIST", test.envArg)
		config, err := LoadConfig(test.args)
		if err != nil {
			if test.expected != "" {
				t.Error(err)
			}
		} else {
			names := make([]string, len(config.Zones))
			for i, z := range config.Zones {
				names[i] = z.Name
			}
			observed := strings.Join(names, ";")
			if observed != test.expected {
				t.Errorf("Expected '%v' to be: '%v'", observed, test.expected)
			}
		}
	}
	os.Setenv("TZ_LIST", oldEnv)
}

func TestSetupZone(t *testing.T) {
	now := time.Now()

	tests := []struct {
		zoneName string
		ok       bool
	}{
		{
			zoneName: "Europe/Paris",
			ok:       true,
		},
		{
			zoneName: "America/London",
			ok:       false,
		},
		{
			// Names should be trimmed in the config, so this should be ok.
			zoneName: " Australia/Sydney",
			ok:       true,
		},
	}
	for _, test := range tests {
		_, err := SetupZone(now, test.zoneName)
		if test.ok != (err == nil) {
			t.Errorf("Expected %v, but got: %v", test.ok, err)
		}
	}
}

func TestSetupZoneWithCustomNames(t *testing.T) {
	now := time.Now()

	tests := []struct {
		zoneName  string
		shortName string
		ok        bool
	}{
		{
			zoneName:  "Europe/Paris,bonjour",
			shortName: "bonjour",
			ok:        true,
		},
		{
			zoneName:  "America/Mexico_City,hola",
			shortName: "hola",
			ok:        true,
		},
		{
			zoneName:  "America/Invalid",
			shortName: "",
			ok:        false,
		},
	}
	for _, test := range tests {
		z, err := SetupZone(now, test.zoneName)
		if test.ok != (err == nil) {
			t.Errorf("Expected %v, but got: %v", test.ok, err)
		}
		if z != nil && !strings.Contains(z.Name, test.shortName) {
			t.Errorf("Expected %v to contain: %v", z.Name, test.shortName)
		}
	}
}
