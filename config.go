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
	"time"
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
	Quit       []string
}

// Config stores app configuration
type Config struct {
	Zones   []*Zone
	Keymaps Keymaps
}

// Function to provide default values for the Config struct
func NewDefaultConfig() Config {
	return Config{
		Zones: DefaultZones,
		Keymaps: Keymaps{
			PrevHour:   []string{"h", "left"},
			NextHour:   []string{"l", "right"},
			PrevDay:    []string{"k", "up"},
			NextDay:    []string{"j", "down"},
			PrevWeek:   []string{"p"},
			NextWeek:   []string{"n"},
			ToggleDate: []string{"d"},
			OpenWeb:    []string{"o"},
			Now:        []string{"t"},
			Quit:       []string{"q", "ctrl+c", "esc"},
		},
	}
}

func LoadConfig(tzConfigs []string) (Config, error) {
	// Apply config file first
	fileConfig, fileError := LoadConfigFile()
	if fileError != nil {
		panic(fileError)
	}

	// Override with env var config
	envConfig, _ := LoadConfigEnv(tzConfigs)

	// Merge configs, with envConfig taking precedence
	mergedConfig := NewDefaultConfig()

	var zones []*Zone

	// Setup with Local time zone
	localZoneName, _ := time.Now().Zone()
	zones = append(zones, &Zone{
		Name:   fmt.Sprintf("(%s) Local", localZoneName),
		DbName: localZoneName,
	})

	// Merge Zones
	if len(envConfig.Zones) > 0 {
		zones = append(zones, envConfig.Zones...)
	} else if len(fileConfig.Zones) > 0 {
		zones = append(zones, fileConfig.Zones...)
	}

	mergedConfig.Zones = zones

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

	return mergedConfig, nil
}
