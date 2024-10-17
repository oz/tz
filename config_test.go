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

func TestConfigKeysDuplicated(t *testing.T) {
	tomlPath := "./config_test_keys_dup.toml"
	_, err := LoadConfig(tomlPath, nil)
	if err == nil {
		t.Errorf("Expected error while reading %s, but didn’t get one", tomlPath)
	}

	expectedError := "Key(s) mapped multiple times in config: q"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected specific error while reading %s, but get a different one: %v", tomlPath, err)
	}
}

func TestLoadConfig(t *testing.T) {
	oldTzList, tzListWasSet := os.LookupEnv("TZ_LIST")
	os.Unsetenv("TZ_LIST")

	tomlPath := "./example-conf.toml"
	_, err := LoadConfig(tomlPath, nil)
	if err != nil {
		t.Errorf("Could not read %s: %v", tomlPath, err)
	}

	if tzListWasSet {
		os.Setenv("TZ_LIST", oldTzList)
		SetupLogger()
	}
}

func TestLoadDefaultConfig(t *testing.T) {
	_, err := LoadDefaultConfig(nil)
		if err != nil {
			t.Fatalf("Could not read default config file: %v", err)
		}
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
		_, err := ReadZoneFromString(now, test.zoneName)
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
		z, err := ReadZoneFromString(now, test.zoneName)
		if test.ok != (err == nil) {
			t.Errorf("Expected %v, but got: %v", test.ok, err)
		}
		if z != nil && !strings.Contains(z.Name, test.shortName) {
			t.Errorf("Expected %v to contain: %v", z.Name, test.shortName)
		}
	}
}
