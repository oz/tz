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
	"slices"
	"strings"
)

// Keymaps represents the key mappings in the TOML file
type Keymaps struct {
	PrevHour   []string
	NextHour   []string
	PrevDay    []string
	NextDay    []string
	PrevWeek   []string
	NextWeek   []string
	ToggleDate []string
	OpenWeb    []string
	Now        []string
	Help       []string
	Quit       []string
}

// Config stores app configuration
type Config struct {
	Zones   []*Zone
	Keymaps Keymaps
}

// Function to provide default values for the Config struct
var DefaultKeymaps = Keymaps{
	PrevHour:   []string{"h", "left"},
	NextHour:   []string{"l", "right"},
	PrevDay:    []string{"H", "shift+left", "pgup", "shift+up", "ctrl+b"},
	NextDay:    []string{"L", "shift+right", "pgdown", "shift+down", "ctrl+f"},
	PrevWeek:   []string{"p", "ctrl+left", "shift+pgup"},
	NextWeek:   []string{"n", "ctrl+right", "shift+pgdown"},
	ToggleDate: []string{"d"},
	OpenWeb:    []string{"o"},
	Now:        []string{"t"},
	Help:       []string{"?"},
	Quit:       []string{"q", "ctrl+c", "esc"},
}

func LoadDefaultConfig(tzConfigs []string) (*Config, error) {
	fileName, fileError := DefaultConfigFile()
	if fileError != nil {
		return nil, fmt.Errorf("File error: %w", fileError)
	}
	return LoadConfig(*fileName, tzConfigs)
}

func LoadConfig(tomlFile string, tzConfigs []string) (*Config, error) {
	// Apply config file first
	fileConfig, fileError := LoadConfigFile(tomlFile)
	if fileError != nil {
		return nil, fmt.Errorf("File error: %w", fileError)
	}

	// Override with env var config
	envConfig, envErr := LoadConfigEnv(tzConfigs)
	if envErr != nil {
		return nil, fmt.Errorf("Env error: %w", envErr)
	}

	// Merge configs, with envConfig taking precedence
	mergedConfig := Config{
		Zones:   []*Zone{DefaultZones[0]},
		Keymaps: DefaultKeymaps,
	}

	// Merge Zones
	var configZones []*Zone
	if len(envConfig.Zones) > 0 {
		configZones = envConfig.Zones
	} else if len(fileConfig.Zones) > 0 {
		configZones = fileConfig.Zones
	} else {
		configZones = DefaultZones[1:]
	}
	mergedConfig.Zones = append(mergedConfig.Zones, configZones...)

	logger.Printf("File zones: %s", fileConfig.Zones)
	logger.Printf("Env zones: %s", envConfig.Zones)
	logger.Printf("Merged zones: %s", mergedConfig.Zones)

	// Merge Keymaps
	if len(fileConfig.Keymaps.PrevHour) > 0 {
		mergedConfig.Keymaps.PrevHour = fileConfig.Keymaps.PrevHour
	}

	if len(fileConfig.Keymaps.NextHour) > 0 {
		mergedConfig.Keymaps.NextHour = fileConfig.Keymaps.NextHour
	}

	if len(fileConfig.Keymaps.PrevDay) > 0 {
		mergedConfig.Keymaps.PrevDay = fileConfig.Keymaps.PrevDay
	}

	if len(fileConfig.Keymaps.NextDay) > 0 {
		mergedConfig.Keymaps.NextDay = fileConfig.Keymaps.NextDay
	}

	if len(fileConfig.Keymaps.PrevWeek) > 0 {
		mergedConfig.Keymaps.PrevWeek = fileConfig.Keymaps.PrevWeek
	}

	if len(fileConfig.Keymaps.NextWeek) > 0 {
		mergedConfig.Keymaps.NextWeek = fileConfig.Keymaps.NextWeek
	}

	if len(fileConfig.Keymaps.ToggleDate) > 0 {
		mergedConfig.Keymaps.ToggleDate = fileConfig.Keymaps.ToggleDate
	}

	if len(fileConfig.Keymaps.OpenWeb) > 0 {
		mergedConfig.Keymaps.OpenWeb = fileConfig.Keymaps.OpenWeb
	}

	if len(fileConfig.Keymaps.Now) > 0 {
		mergedConfig.Keymaps.Now = fileConfig.Keymaps.Now
	}

	if len(fileConfig.Keymaps.Quit) > 0 {
		mergedConfig.Keymaps.Quit = fileConfig.Keymaps.Quit
	}

	allKeymaps := [][]string {
		mergedConfig.Keymaps.PrevHour,
		mergedConfig.Keymaps.NextHour,
		mergedConfig.Keymaps.PrevDay,
		mergedConfig.Keymaps.NextDay,
		mergedConfig.Keymaps.PrevWeek,
		mergedConfig.Keymaps.NextWeek,
		mergedConfig.Keymaps.ToggleDate,
		mergedConfig.Keymaps.OpenWeb,
		mergedConfig.Keymaps.Now,
		mergedConfig.Keymaps.Help,
		mergedConfig.Keymaps.Quit,
	}
	var keysUsed = make(map[string]bool)
	var keysDuplicated []string
	for _, keys := range allKeymaps {
		for _, key := range keys {
			if _, used := keysUsed[key]; used == false {
				keysUsed[key] = true
			} else {
				keysDuplicated = append(keysDuplicated, key)
			}
		}
	}
	if len(keysDuplicated) > 0 {
		slices.Sort(keysDuplicated)
		return nil, fmt.Errorf("Key(s) mapped multiple times in config: %v", strings.Join(keysDuplicated, " "))
	}

	return &mergedConfig, nil
}
