package main

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/tkuchiki/go-timezone"
)

//go:embed rawData.json.gz
var rawData []byte

type TZCountry struct {
	TzName             string   `json:"name"`
	AlternativeName    string   `json:"alternativeName"`
	Group              []string `json:"group"`
	ContinentCode      string   `json:"continentCode"`
	ContinentName      string   `json:"continentName"`
	CountryName        string   `json:"countryName"`
	CountryCode        string   `json:"countryCode"`
	MainCities         []string `json:"mainCities"`
	RawOffsetInMinutes int      `json:"rawOffsetInMinutes"`
	Abbreviation       string   `json:"abbreviation"`
	RawFormat          string   `json:"rawFormat"`
}

func loadCountries() map[string]TZCountry {
	b := bytes.NewReader(rawData)
	gz, _ := gzip.NewReader(b)
	defer gz.Close()
	var countries []TZCountry
	json.NewDecoder(gz).Decode(&countries)
	countryMap := make(map[string]TZCountry)
	for _, c := range countries {
		countryMap[c.TzName] = c
	}
	return countryMap
}

func doTZList() {
	tz := timezone.New()
	var names []string
	for name, info := range tz.TzInfos() {
		if info.IsDeprecated() {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	countries := loadCountries()

	for _, name := range names {
		info, _ := tz.GetTzInfo(name)
		out := fmt.Sprintf("%9s (%s) :: %s", info.ShortStandard(), info.StandardOffsetHHMM(), name)
		c := countries[name]
		if c.CountryName != "" {
			out += fmt.Sprintf(" (%s)", c.CountryName)
		} else {
			// search group for matching tz
			for _, c := range countries {
				found := false
				for _, t := range c.Group {
					if t == name {
						found = true
						out += fmt.Sprintf(" (%s)", c.CountryName)
						break
					}
				}
				if found {
					break
				}
			}
		}
		fmt.Println(out)
	}
}
