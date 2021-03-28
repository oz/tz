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
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Config stores app configuration
type Config struct {
	Zones     []*Zone
	Time24    bool
	ExitQuick bool
}

// Usage text
const usage = `Usage:
    tz [-l <local name>] [-24] [-q]

Options:
    -l NAME            Replace "Local" timezone name with "NAME"
    -q                 Show timezones and exit immediately
    -24                Display times in 24h time format
    -h                 Show this help text`

// LoadConfig from environment
func LoadConfig() (*Config, error) {
	flag.Usage = func() { _, _ = fmt.Fprintf(os.Stderr, "%s\n", usage) }
	conf := Config{
		Zones:  DefaultZones,
		Time24: false,
	}

	// Parse flags
	localIdentifier := flag.String("l", "Local", `Override "Local" with "NAME"`)
	flag.BoolVar(&conf.Time24, "24", false, `Display times in 24h format`)
	flag.BoolVar(&conf.ExitQuick, "q", false, "exit immediately")
	flag.Parse()

	tzList := os.Getenv("TZ_LIST")
	if tzList == "" {
		conf.Zones[0].Name = *localIdentifier
		return &conf, nil
	}
	tzConfigs := strings.Split(tzList, ",")
	if len(tzConfigs) == 0 {
		return &conf, nil
	}
	zones := make([]*Zone, len(tzConfigs)+1)

	// Setup with Local time zone
	now := time.Now()
	localZoneName, offset := now.Zone()
	zones[0] = &Zone{
		Name:   fmt.Sprintf("(%s) %s", localZoneName, *localIdentifier),
		DbName: localZoneName,
		Offset: offset / 3600,
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
	names := strings.Split(zoneConf, ";")
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
	then := now.In(loc)
	shortName, offset := then.Zone()
	return &Zone{
		DbName: loc.String(),
		Name:   fmt.Sprintf("(%s) %s", shortName, name),
		Offset: offset / 3600,
	}, nil
}
