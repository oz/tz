package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the entire TOML configuration
type ConfigFile struct {
	Header  string            `toml:"header"`
	Zones   []ConfigFileZone  `toml:"zones"`
	Keymaps ConfigFileKeymaps `toml:"keymaps"`
}

// Zone represents a single zone entry in the TOML file
type ConfigFileZone struct {
	ID   string `toml:"id"`
	Name string `toml:"name"`
}

// Keymaps represents the key mappings in the TOML file
type ConfigFileKeymaps struct {
	PrevHour   []string `toml:"prev_hour"`
	NextHour   []string `toml:"next_hour"`
	PrevDay    []string `toml:"prev_day"`
	NextDay    []string `toml:"next_day"`
	PrevWeek   []string `toml:"prev_week"`
	NextWeek   []string `toml:"next_week"`
	ToggleDate []string `toml:"toggle_date"`
	OpenWeb    []string `toml:"open_web"`
	Now        []string `toml:"now"`
}

func setupZone(now time.Time, zoneConf ConfigFileZone) (*Zone, error) {
	name := zoneConf.Name
	dbName := zoneConf.ID

	loc, err := time.LoadLocation(dbName)
	if err != nil {
		return nil, fmt.Errorf("looking up zone %s: %w", dbName, err)
	}
	if name == "" {
		name = loc.String()
	}
	then := now.In(loc)
	shortName, _ := then.Zone()
	return &Zone{
		DbName: loc.String(),
		Name:   fmt.Sprintf("(%s) %s", shortName, name),
	}, nil
}

func LoadConfigFile() (*Config, error) {
	conf := Config{
		Zones: DefaultZones,
	}

	// Expand the ~ to the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configFilePath := filepath.Join(homeDir, ".config", "tz", "conf.toml")

	// Read the TOML file
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}

	// Unmarshal the TOML data into the Config struct
	var config ConfigFile
	err = toml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	zones := make([]*Zone, len(config.Zones)+1)

	// Setup with Local time zone
	localZoneName, _ := time.Now().Zone()
	zones[0] = &Zone{
		Name:   fmt.Sprintf("(%s) Local", localZoneName),
		DbName: localZoneName,
	}

	// Add zones from config file
	for i, zoneConf := range config.Zones {
		zone, err := setupZone(time.Now(), zoneConf)
		if err != nil {
			return nil, err
		}
		zones[i+1] = zone
	}

	conf.Zones = zones
	conf.Keymaps = Keymaps(config.Keymaps)

	return &conf, nil
}
