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
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdateIncHour(t *testing.T) {
	// "l" key -> go right
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'l'},
		Alt:   false,
	}

	tests := []struct {
		startHour int
		nextHour  int
	}{
		{startHour: 0, nextHour: 1},
		{startHour: 1, nextHour: 2},
		// ...
		{startHour: 23, nextHour: 0},
	}

	for _, test := range tests {
		m := model{
			zones: DefaultZones,
			hour:  test.startHour,
		}
		nextState, cmd := m.Update(msg)
		if cmd != nil {
			t.Errorf("Expected nil Cmd, but got %v", cmd)
			return
		}
		h := nextState.(model).hour
		if h != test.nextHour {
			t.Errorf("Expected %d, but got %d", test.nextHour, h)
		}
	}
}

func TestUpdateDecHour(t *testing.T) {
	// "h" key -> go right
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'h'},
		Alt:   false,
	}

	tests := []struct {
		startHour int
		nextHour  int
	}{
		{startHour: 23, nextHour: 22},
		{startHour: 22, nextHour: 21},
		// ...
		{startHour: 0, nextHour: 23},
	}

	for _, test := range tests {
		m := model{
			zones: DefaultZones,
			hour:  test.startHour,
		}
		nextState, cmd := m.Update(msg)
		if cmd != nil {
			t.Errorf("Expected nil Cmd, but got %v", cmd)
			return
		}
		h := nextState.(model).hour
		if h != test.nextHour {
			t.Errorf("Expected %d, but got %d", test.nextHour, h)
		}
	}
}

func TestUpdateQuitMsg(t *testing.T) {
	// "q" key -> quit
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
		Alt:   false,
	}

	m := model{
		zones: DefaultZones,
		hour:  10,
	}
	_, cmd := m.Update(msg)
	if cmd == nil {
		t.Errorf("Expected tea.Quit Cmd, but got %v", cmd)
		return
	}
	// tea.Quit is a function, we can't really test with == here, and
	// calling it is getting into internal territory.
}
