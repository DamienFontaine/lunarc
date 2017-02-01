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

package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus/hooks/test"
)

func TestSingleFileNormal(t *testing.T) {
	request, _ := http.NewRequest("GET", "robot.txt", nil)

	w := httptest.NewRecorder()
	SingleFile("robot.txt").ServeHTTP(w, request)

	if w.Code != http.StatusOK {
		t.Fatalf("Non expected code: %v", w.Code)
	}
}

func TestLoggingNormal(t *testing.T) {
	request, _ := http.NewRequest("GET", "robot.txt", nil)

	logger, hook := test.NewNullLogger()

	next := SingleFile("robot.txt")

	w := httptest.NewRecorder()
	Logging(next, logger).ServeHTTP(w, request)

	if len(hook.Entries) != 1 {
		t.Fatalf("Must return 1 but : %v", len(hook.Entries))
	}
}

func TestSingleFileNotFound(t *testing.T) {
	request, _ := http.NewRequest("GET", "robot.txt", nil)

	w := httptest.NewRecorder()
	SingleFile("robots.txt").ServeHTTP(w, request)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Non expected code: %v", w.Code)
	}
}
