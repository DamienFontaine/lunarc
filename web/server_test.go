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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

func GetLoggingHTTPServer(t *testing.T, env string) (s *Server) {
	s, err := NewServer("config.yml", env)
	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
	m := s.Handler.(*LoggingServeMux)
	m.Handle("/", SingleFile("hello.html"))
	return
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

func TestNewServerWithNoLog(t *testing.T) {
	server := getHTTPServer(t, "testNoLog")

	if server.Config.Port != 8888 {
		t.Fatalf("Non expected server port: %v != %v", 8888, server.Config.Port)
	}
	if server.Config.Jwt.Key != "LunarcSecretKey" {
		t.Fatalf("Non expected server Jwt secret key: %v != %v", "LunarcSecretKey", server.Config.Jwt.Key)
	}
}

func TestNewLoggingServeMux(t *testing.T) {
	server := GetLoggingHTTPServer(t, "test")

	hook := test.NewGlobal()
	go server.Start()
	time.Sleep(time.Second * 3)
	errs := make(chan error, 1)

	go func() {
		_, err := http.Get("http://localhost:8888/")

		time.Sleep(time.Second * 1)

		if err != nil {
			errs <- err
			return
		}
		go server.Stop()
	}()

	select {
	case <-server.Done:
		if len(hook.Entries) != 3 {
			for _, entry := range hook.Entries {
				log.Printf("Entry: %v", entry)
			}
			t.Fatalf("Must return 3 but : %d", len(hook.Entries))
		}
	case err := <-errs:
		go server.Stop()
		<-server.Done
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	case err := <-server.Error:
		go server.Stop()
		<-server.Done
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}
}

func TestStart(t *testing.T) {
	server := getHTTPServer(t, "test")

	go server.Start()

	time.Sleep(time.Second * 3)

	errs := make(chan error, 1)

	go func() {
		resp, err := http.Get("http://localhost:8888/")
		if err != nil {
			errs <- err
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errs <- err
			return
		}

		if !strings.Contains(string(body), "Lunarc") {
			errs <- fmt.Errorf("Body must contain Lunarc word but : %v", body)
			return
		}
		go server.Stop()
	}()

	select {
	case <-server.Done:
	case err := <-errs:
		go server.Stop()
		<-server.Done
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	case err := <-server.Error:
		go server.Stop()
		<-server.Done
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}
}

func TestStartWithSSLNormal(t *testing.T) {
	server := getHTTPServer(t, "ssl")

	go server.Start()

	time.Sleep(time.Second * 3)

	errs := make(chan error, 1)

	go func() {
		client := &http.Client{Transport: &http2.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true, NextProtos: []string{"h2"}}}}

		response, err := client.Get("https://localhost:8888/")
		if err != nil {
			errs <- err
			return
		}

		if response.TLS == nil {
			errs <- errors.New("This connection must be in HTTPS")
			return
		}

		if strings.Compare(response.Proto, "HTTP/2.0") != 0 {
			errs <- fmt.Errorf("Must use HTTP/2 but use : %v", response.Proto)
			return
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			errs <- err
			return
		}

		if !strings.Contains(string(body), "Lunarc") {
			errs <- fmt.Errorf("Body must contain Lunarc word but : %v", body)
			return
		}

		go server.Stop()
	}()

	select {
	case <-server.Done:
	case err := <-errs:
		go server.Stop()
		<-server.Done
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	case err := <-server.Error:
		go server.Stop()
		<-server.Done
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}
}

func TestStopNormal(t *testing.T) {
	server := getHTTPServer(t, "test")

	go server.Start()

	time.Sleep(time.Second * 3)

	errs := make(chan error, 1)

	go func() {
		resp, err := http.Get("http://localhost:8888/")
		if err != nil {
			errs <- err
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errs <- err
			return
		}

		if !strings.Contains(string(body), "Lunarc") {
			errs <- fmt.Errorf("Body must contain Lunarc word but : %v", body)
			return
		}
		go server.Stop()
	}()

	select {
	case <-server.Done:
	case err := <-errs:
		go server.Stop()
		<-server.Done
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	case err := <-server.Error:
		go server.Stop()
		<-server.Done
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	resp, err := http.Get("http://localhost:8888/")
	if err == nil {
		t.Fatalf("Error expected: Not Found: %v", resp)
	}

	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		t.Fatalf("Error : %v", err)
	}
	defer l.Close()
}

func TestStartWithSSLNoCertError(t *testing.T) {
	server := getHTTPServer(t, "nocertssl")

	go server.Start()

	err := <-server.Error
	if err == nil {
		t.Fatalf("Expected error: Le fichier spécifié est introuvable")
	}
	<-server.Done
}

func TestStartWithSSLNoKeyError(t *testing.T) {
	server := getHTTPServer(t, "nokeyssl")

	go server.Start()

	err := <-server.Error
	if err == nil {
		t.Fatalf("Expected error: Le fichier spécifié est introuvable")
	}
	<-server.Done
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
	}
}
