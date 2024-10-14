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
	"time"
)

// Config stores app configuration
type Config struct {
	Zones []*Zone
}

// LoadConfig from environment
func LoadConfig(tzConfigs []string) (*Config, error) {
	conf := Config{
		Zones: DefaultZones,
	}

	if len(tzConfigs) == 0 {
		tzList := os.Getenv("TZ_LIST")
		if tzList == "" {
			return &conf, nil
		}
		tzConfigs = strings.Split(tzList, ";")
		if len(tzConfigs) == 0 {
			return &conf, nil
		}
	}
	zones := make([]*Zone, len(tzConfigs)+1)

	// Setup with Local time zone
	now := time.Now()
	localZoneName, _ := now.Zone()
	zones[0] = &Zone{
		Loc:    now.Location(),
		Name:   "Local",
		DbName: localZoneName,
	}

	// Add zones from TZ_LIST
	for i, zoneConf := range tzConfigs {
		zone, err := SetupZone(now, zoneConf)
		if err != nil {
			return nil, err
		}
		zones[i+1] = zone
	}
	conf.Zones = zones

	return &conf, nil
}

// SetupZone from current time and a zoneConf string
func SetupZone(now time.Time, zoneConf string) (*Zone, error) {
	names := strings.Split(zoneConf, ",")
	dbName := strings.Trim(names[0], " ")
	var name string
	if len(names) == 2 {
		name = names[1]
	}

	loc, err := time.LoadLocation(dbName)
	if err != nil {
		return nil, fmt.Errorf("looking up zone %s: %w", dbName, err)
	}
	if name == "" {
		name = loc.String()
	}
	return &Zone{
		Loc:    loc,
		DbName: loc.String(),
		Name:   name,
	}, nil
}
