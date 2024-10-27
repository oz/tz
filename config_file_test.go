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
	"testing"
	"time"
)

func TestDefaultConfigFile(t *testing.T) {
	_, err := DefaultConfigFile()
	if err != nil {
		t.Fatalf("Could not read default config file: %v", err)
	}
}

func TestEmptyConfigFile(t *testing.T) {
	now := time.Now()
	tomlPath := "./testdata/config_file/empty.toml"
	_, err := LoadConfigFile(tomlPath, now)
	if err != nil {
		t.Fatalf("Unexpected error reading empty config %s: %v", tomlPath, err)
	}
}

func TestExampleConfigFile(t *testing.T) {
	now := time.Now()
	tomlPath := "./example-conf.toml"
	config, err := LoadConfigFile(tomlPath, now)
	if err != nil {
		t.Fatalf("Could not read test config from %s: %v", tomlPath, err)
	}

	if len(config.Zones) < 4 {
		t.Errorf("Expected at least 4 zones in %s, found %v", tomlPath, len(config.Zones))
	}

	if len(config.Keymaps.OpenWeb) < 2 {
		t.Errorf("Expected at least 2 keys for open_web in %s, found %v", tomlPath, len(config.Keymaps.OpenWeb))
	}
}

func TestInvalidConfigFile(t *testing.T) {
	now := time.Now()
	tomlPath := "./testdata/config_file/invalid.toml"

	{
		_, err := LoadConfigFile(tomlPath, now)
		failUnlessExpectedError(t, err, "toml: expected character ]", "reading %s", tomlPath)
	}

	{
		_, err := LoadConfig(tomlPath, nil)
		failUnlessExpectedError(t, err, "toml: expected character ]", "reading %s", tomlPath)
	}
}

func TestUnknownZoneConfigFile(t *testing.T) {
	now := time.Now()
	tomlPath := "./testdata/config_file/unknown_zone.toml"
	_, err := LoadConfigFile(tomlPath, now)
	failUnlessExpectedError(t, err, "unknown time zone !", "reading %s", tomlPath)
}
