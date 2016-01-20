// Copyright (c) - Damien Fontaine <damien.fontaine@lineolia.net>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>

package utils

import (
	"testing"
)

func TestSanitizeAccent(t *testing.T) {
	s := "ŒOEœoeàa"
	r := SanitizeAccent(s)

	if r != "OEOEoeoeaa" {
		t.Fatalf("Non expected string %v", r)
	}
}

func TestSanitizeTitle(t *testing.T) {
	s := " Test  d'un titre être ou ne pas été !"
	p := SanitizeTitle(s)

	if p == "" {
		t.Fatalf("Non expected title %v", p)
	}
	if p != "test-d-un-titre-etre-ou-ne-pas-ete" {
		t.Fatalf("Non expected title %v", p)
	}
}
