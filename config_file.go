package main

import (
	"fmt"
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
	PrevMinute []string `toml:"prev_minute"`
	NextMinute []string `toml:"next_minute"`
	ZeroMinute []string `toml:"zero_minute"`
	PrevHour   []string `toml:"prev_hour"`
	NextHour   []string `toml:"next_hour"`
	PrevDay    []string `toml:"prev_day"`
	NextDay    []string `toml:"next_day"`
	PrevWeek   []string `toml:"prev_week"`
	NextWeek   []string `toml:"next_week"`
	PrevLine   []string `toml:"prev_line_select"`
	NextLine   []string `toml:"next_line_select"`
	PrevFStyle []string `toml:"prev_format_style"`
	NextFStyle []string `toml:"next_format_style"`
	PrevZStyle []string `toml:"prev_zone_style"`
	NextZStyle []string `toml:"next_zone_style"`
	ToggleDate []string `toml:"toggle_date"`
	OpenWeb    []string `toml:"open_web"`
	Now        []string `toml:"now"`
	Help       []string `toml:"help"`
	Quit       []string `toml:"quit"`
}

func ReadZonesFromFile(now time.Time, zoneConf ConfigFileZone) (*Zone, error) {
	name := zoneConf.Name
	dbName := zoneConf.ID

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

func DefaultConfigFile() (*string, error) {
	// Return early if we can't find a home dir.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configFilePath := filepath.Join(homeDir, ".config", "tz", "conf.toml")
	return &configFilePath, nil
}

func LoadConfigFile(configFilePath string, now time.Time) (*Config, error) {
	conf := Config{}

	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		// Ignore unreadable config file.
		logger.Printf("Config file '%s' not found. Skipping...\n", configFilePath)
		return &conf, nil
	}

	var config ConfigFile
	if err = toml.Unmarshal(configFile, &config); err != nil {
		return nil, fmt.Errorf("Parsing %s: %w\n", configFilePath, err)
	}

	// Add zones from config file
	zones := make([]*Zone, len(config.Zones))
	for i, zoneConf := range config.Zones {
		zone, err := ReadZonesFromFile(now, zoneConf)
		if err != nil {
			return nil, err
		}
		zones[i] = zone
	}

	conf.Zones = zones
	conf.Keymaps = Keymaps(config.Keymaps)

	return &conf, nil
}
