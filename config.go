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
	Zones []Zone
}

// LoadConfig from os.UserConfigDir
func LoadConfig() (*Config, error) {
	conf := Config{
		Zones: DefaultZones,
	}

	zoneEnv := os.Getenv("TZ_LIST")
	if zoneEnv == "" {
		return &conf, nil
	}
	zoneNames := strings.Split(zoneEnv, ",")
	if len(zoneNames) == 0 {
		return &conf, nil
	}
	zones := make([]Zone, len(zoneNames)+1)

	now := time.Now()
	localZoneName, offset := now.Zone()

	zones[0] = Zone{
		Name:   fmt.Sprintf("(%s) Local", localZoneName),
		DbName: localZoneName,
		Offset: offset / 3600,
	}
	for i, name := range zoneNames {
		loc, err := time.LoadLocation(name)
		if err != nil {
			return nil, fmt.Errorf("looking up zone %s: %w", name, err)
		}
		then := now.In(loc)
		shortName, offset := then.Zone()
		zones[i+1] = Zone{
			DbName: loc.String(),
			Name:   fmt.Sprintf("(%s) %s", shortName, loc),
			Offset: offset / 3600,
		}
	}
	conf.Zones = zones
	return &conf, nil
}
