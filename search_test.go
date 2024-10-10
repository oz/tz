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
	"strings"
	"testing"
)

func TestSearchAndPrintFilter(t *testing.T) {
	var builder strings.Builder
	SearchZones(strings.ToLower("UTC")).Print(&builder)
	observed := builder.String()
	expected := "  UTC (+00:00) :: UTC\n"
	if !strings.Contains(observed, expected) {
		t.Errorf(fmt.Sprintf("Missing '%v' in '%v'", strings.TrimSpace(expected), strings.TrimSpace(observed)))
	}
}
