package main

import (
	"fmt"
	"sort"

	"github.com/tkuchiki/go-timezone"
)

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
	for _, name := range names {
		info, _ := tz.GetTzInfo(name)
		fmt.Printf("%9s (%s) :: %s\n", info.ShortStandard(), info.StandardOffsetHHMM(), name)
	}
}
