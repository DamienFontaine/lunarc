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
	"net"
	"net/http"
	"testing"
	"time"
)

func TestInitialize(t *testing.T) {
	server, err := NewWebServer("config.yml", "test")
	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
	if server.GetContext().GetCnf().Server.Port != 8888 {
		t.Fatalf("Non expected server port: %v != %v", 8888, server.GetContext().GetCnf().Server.Port)
	}
	if server.GetContext().GetCnf().Server.Jwt.Key != "LunarcSecretKey" {
		t.Fatalf("Non expected server Jwt secret key: %v != %v", "LunarcSecretKey", server.GetContext().GetCnf().Server.Port)
	}
}

func TestStart(t *testing.T) {
	server, err := NewWebServer("config.yml", "test")
	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
	go server.Start()

	time.Sleep(time.Second * 3)

	_, err = http.Get("http://localhost:8888/")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	server.Stop()
}

func TestStartWithError(t *testing.T) {
	l, err := net.Listen("tcp", ":8888")
	defer l.Close()
	go http.Serve(l, nil)

	server, _ := NewWebServer("config.yml", "test")

	go server.Start()

	err = <-server.err

	if err == nil {
		t.Fatalf("Expected error: listen tcp :8888: bind: address already in use")
	}
}

func TestStopNormal(t *testing.T) {
	server, err := NewWebServer("config.yml", "test")
	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
	go server.Start()

	time.Sleep(time.Second * 3)

	_, err = http.Get("http://localhost:8888/")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	server.Stop()

	time.Sleep(time.Second * 3)

	resp, err := http.Get("http://localhost:8888/")
	if err == nil {
		t.Fatalf("Error expected: Not Found: %v", resp)
	}
}

func TestStopUnstarted(t *testing.T) {
	server, err := NewWebServer("config.yml", "test")
	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}

	server.Stop()
}
