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
)

func TestDefaultConfigFile(t *testing.T) {
	_, err := DefaultConfigFile()
	if err != nil {
		t.Fatalf("Could not read default config file: %v", err)
	}
}

func TestExampleConfigFile(t *testing.T) {
	tomlPath := "./example-conf.toml"
	config, err := LoadConfigFile(tomlPath)
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
