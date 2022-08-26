package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/tkuchiki/go-timezone"
)

// Zone search results
type ZoneSearchResults map[string]*timezone.TzInfo

// List of zones sorted alphabetically.
func (zsr ZoneSearchResults) SortedNames() []string {
	sorted := make([]string, 0, len(zsr))
	for name := range zsr {
		sorted = append(sorted, name)
	}
	sort.Strings(sorted)

	return sorted
}

// Print formatted ZoneSearchResults to the chosen Writer ; typically
// os.Stdout.
func (zsr ZoneSearchResults) Print(w io.Writer) {
	sorted := zsr.SortedNames()
	for i := range sorted {
		name := sorted[i]
		ti := zsr[name]
		fmt.Fprintf(w, "%5s (%s) :: %s\n",
			ti.ShortStandard(),
			ti.StandardOffsetHHMM(),
			name)
	}
}

// Find zones matching a query. An empty query string returns all zones.
func SearchZones(q string) ZoneSearchResults {
	// TODO Each call to timezone.New() allocs a fresh list of timezones:
	//      for now, avoid calling SearchZones too much.
	t := timezone.New()
	filter := q != ""
	matches := map[string]*timezone.TzInfo{}

	for abbr, zones := range t.Timezones() {
		for _, name := range zones {
			if filter &&
				!strings.Contains(strings.ToLower(name), q) &&
				!strings.Contains(strings.ToLower(abbr), q) {
				continue
			}

			ti, err := t.GetTzInfo(name)
			// That should not happen too often.
			if err != nil {
				panic(err)
			}

			matches[name] = ti
		}
	}
	return matches
}
