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
}

// Config stores app configuration
type Config struct {
	Zones   []*Zone
	Keymaps Keymaps
}
