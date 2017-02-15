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
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Sirupsen/logrus/hooks/test"

	"golang.org/x/net/http2"
)

func getHTTPServer(t *testing.T, env string) (s *Server) {
	s, err := NewServer("config.yml", env)
	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
	m := s.Handler.(*LoggingServeMux)
	m.Handle("/", SingleFile("hello.html"))
	return
}

func TestNewServerWithNoLog(t *testing.T) {
	server := getHTTPServer(t, "testNoLog")

	if server.Config.Port != 8888 {
		t.Fatalf("Non expected server port: %v != %v", 8888, server.Config.Port)
	}
	if server.Config.Jwt.Key != "LunarcSecretKey" {
		t.Fatalf("Non expected server Jwt secret key: %v != %v", "LunarcSecretKey", server.Config.Jwt.Key)
	}
}

func TestNewServer(t *testing.T) {
	server := getHTTPServer(t, "test")

	if server.Config.Port != 8888 {
		t.Fatalf("Non expected server port: %v != %v", 8888, server.Config.Port)
	}
	if server.Config.Jwt.Key != "LunarcSecretKey" {
		t.Fatalf("Non expected server Jwt secret key: %v != %v", "LunarcSecretKey", server.Config.Jwt.Key)
	}
}

func TestStart(t *testing.T) {
	server := getHTTPServer(t, "test")

	go server.Start()

	time.Sleep(time.Second * 3)

	resp, err := http.Get("http://localhost:8888/")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(string(body), "Lunarc") {
		t.Fatalf("Body must contain Lunarc word but : %v", body)
	}

	go server.Stop()
	<-server.Done
}

func TestStartWithSSLNormal(t *testing.T) {
	server := getHTTPServer(t, "ssl")

	go server.Start()

	time.Sleep(time.Second * 3)

	client := &http.Client{Transport: &http2.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true, NextProtos: []string{"h2"}}}}

	response, err := client.Get("https://localhost:8888/")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.TLS == nil {
		t.Fatalf("This connection must be in HTTPS")
	}

	if strings.Compare(response.Proto, "HTTP/2.0") != 0 {
		t.Fatalf("Must use HTTP/2 but use : %v", response.Proto)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(string(body), "Lunarc") {
		t.Fatalf("Body must contain Lunarc word but : %v", body)
	}

	go server.Stop()
	<-server.Done
}

func TestStartWithSSLNoCertError(t *testing.T) {
	server := getHTTPServer(t, "nocertssl")

	go server.Start()

	err := <-server.Error
	if err == nil {
		t.Fatalf("Expected error: Le fichier spécifié est introuvable")
	}
}

func TestStartWithSSLNoKeyError(t *testing.T) {
	server := getHTTPServer(t, "nokeyssl")

	go server.Start()

	err := <-server.Error
	if err == nil {
		t.Fatalf("Expected error: Le fichier spécifié est introuvable")
	}
}

func TestStartWithError(t *testing.T) {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		t.Fatalf("Error during test preparation : %v", err)
	}

	done := make(chan struct{}, 1)

	go func() {
		err = http.Serve(l, nil)
		if err != nil {
			close(done)
		}
	}()

	server := getHTTPServer(t, "test")

	go server.Start()

	err = <-server.Error

	if err == nil {
		t.Fatalf("Expected error: listen tcp :8888: bind: address already in use")
	}
	l.Close()
	<-done
}

func TestStartWithSSLAndError(t *testing.T) {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		t.Fatalf("Error during test preparation : %v", err)
	}

	done := make(chan struct{}, 1)

	go func() {
		err = http.Serve(l, nil)
		if err != nil {
			close(done)
		}
	}()

	server := getHTTPServer(t, "ssl")

	go server.Start()

	err = <-server.Error

	if err == nil {
		t.Fatalf("Expected error: listen tcp :8888: bind: address already in use")
	}

	l.Close()
	<-done
}

func TestStopNormal(t *testing.T) {
	server := getHTTPServer(t, "test")

	go server.Start()

	time.Sleep(time.Second * 2)

	resp, err := http.Get("http://localhost:8888/")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(string(body), "Lunarc") {
		t.Fatalf("Body must contain Lunarc word but : %v", body)
	}

	go server.Stop()
	<-server.Done

	resp, err = http.Get("http://localhost:8888/")
	if err == nil {
		t.Fatalf("Error expected: Not Found: %v", resp)
	}

	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		t.Fatalf("Error : %v", err)
	}
	defer l.Close()
	go http.Serve(l, nil)
}

func TestStopUnstarted(t *testing.T) {
	server := getHTTPServer(t, "test")

	go server.Stop()

	select {
	case result := <-server.Done:
		if result {
			t.Fatalf("Non expected behavior")
		}
	case err := <-server.Error:
		if err == nil {
			t.Fatalf("Non expected error")
		}
		return
	}
}

func GetLoggingHTTPServer(t *testing.T, env string) (s *Server) {
	s, err := NewServer("config.yml", env)
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
