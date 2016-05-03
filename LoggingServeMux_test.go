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

package lunarc

import (
	"net/http"
	"testing"

	"github.com/Sirupsen/logrus/hooks/test"
)

func GetLoggingHTTPServer(t *testing.T, env string) (s *WebServer) {
	s, err := NewWebServer("config.yml", env)
	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
	return
}

func TestNewLoggingServeMux(t *testing.T) {
	server := GetLoggingHTTPServer(t, "test")

	go server.Start()

	hook := test.NewGlobal()

	_, err := http.Get("http://localhost:8888/")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(hook.Entries) != 1 {
		t.Fatalf("Must return 1 but : %d", len(hook.Entries))
	}

	go server.Stop()
	<-server.Done
}
