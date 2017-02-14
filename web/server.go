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
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"

	"golang.org/x/net/http2"

	"github.com/Sirupsen/logrus"
	log "github.com/Sirupsen/logrus"
)

const logFilename = "lunarc.log"

//IServer is an http.ServeMux with a Context.
type IServer interface {
	Start() error
	Stop()
}

//Server is a Server with a specialize Context.
type Server struct {
	http.Server
	Config      Config
	Error       chan error
	Done        chan bool
	quit        chan bool
	interrupt   chan os.Signal
	connections map[net.Conn]http.ConnState
}

//NewServer create a new instance of Server
func NewServer(filename string, environment string) (server *Server, err error) {
	conf, err := GetConfig(filename, environment)

	logFile, err := os.OpenFile(conf.Log.File+logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Warningf("Can't open logfile: %v", err)
	} else {
		log.SetOutput(logFile)
	}

	level := log.ErrorLevel
	if strings.Compare(conf.Log.Level, "") != 0 {
		level, _ = log.ParseLevel(conf.Log.Level)
	} else {
		log.Infof("Log Level: %v", level)
	}
	log.SetLevel(level)

	server = &Server{Config: conf, Done: make(chan bool, 1), Error: make(chan error, 1), Server: http.Server{Handler: NewLoggingServeMux(conf)}, quit: make(chan bool)}
	return
}

//Start the server.
func (s *Server) Start() (err error) {
	log.Infof("Lunarc is starting on port :%d", s.Config.Port)
	go func() {
		var l net.Listener

		l, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Config.Port))
		if err != nil {
			log.Errorf("Error: %v", err)
			s.Error <- err
			return
		}

		// Track connection state
		add := make(chan net.Conn)
		active := make(chan net.Conn)
		idle := make(chan net.Conn)
		remove := make(chan net.Conn)

		s.ConnState = func(conn net.Conn, state http.ConnState) {
			switch state {
			case http.StateNew:
				add <- conn
			case http.StateActive:
				active <- conn
			case http.StateIdle:
				idle <- conn
			case http.StateClosed, http.StateHijacked:
				remove <- conn
			}
		}

		shutdown := make(chan chan struct{})
		go s.handleConnections(add, active, idle, remove, shutdown)
		go s.handleInterrupt(l, shutdown)

		if len(s.Config.SSL.Certificate) > 0 && len(s.Config.SSL.Key) > 0 {
			config := tls.Config{}

			config.Certificates = make([]tls.Certificate, 1)
			config.Certificates[0], err = tls.LoadX509KeyPair(s.Config.SSL.Certificate, s.Config.SSL.Key)
			if err != nil {
				log.Errorf("%v", err)
				l.Close()
				s.Error <- err
				return
			}

			s.TLSConfig = &config

			err = http2.ConfigureServer(&s.Server, nil)
			if err != nil {
				log.Errorf("%v", err)
				l.Close()
				s.Error <- err
				return
			}
			err = s.Serve(tls.NewListener(l, &config))
			if err != nil {
				close(s.quit)
				return
			}
		} else {
			err = s.Serve(l)
			if err != nil {
				close(s.quit)
				return
			}
		}
	}()

	<-s.quit
	s.interrupt <- syscall.SIGINT
	return
}

func (s *Server) handleConnections(add, active, idle, remove chan net.Conn, shutdown chan chan struct{}) {
	var done chan struct{}
	s.connections = map[net.Conn]http.ConnState{}
	for {
		select {
		case conn := <-add:
			s.connections[conn] = http.StateNew
		case conn := <-active:
			s.connections[conn] = http.StateActive
		case conn := <-remove:
			delete(s.connections, conn)
			if done != nil && len(s.connections) == 0 {
				done <- struct{}{}
				return
			}
		case conn := <-idle:
			s.connections[conn] = http.StateIdle
			if done != nil {
				conn.Close()
				if len(s.connections) == 0 {
					done <- struct{}{}
				}
			}
		case done = <-shutdown:
			for k, v := range s.connections {
				if v == http.StateIdle {
					k.Close()
					delete(s.connections, k)
				}
			}
			if len(s.connections) == 0 {
				done <- struct{}{}
			}
			return
		}
	}
}

func (s *Server) handleInterrupt(listener net.Listener, shutdown chan chan struct{}) {
	if s.interrupt == nil {
		s.interrupt = make(chan os.Signal, 1)
	}
	<-s.interrupt
	s.SetKeepAlivesEnabled(false)
	err := listener.Close()
	if err != nil {
		log.Errorf("Error on listener close: %v", err)
	}
	<-s.quit
	done := make(chan struct{})
	shutdown <- done
	<-done
	listener = nil
	s.interrupt = nil
	log.Info("Lunarc terminated.")
	s.Done <- true
}

//Stop the server.
func (s *Server) Stop() {
	if s.interrupt != nil && s.quit != nil {
		log.Info("Lunarc is stopping...")
		s.quit <- true
	} else {
		log.Info("Lunarc is not running")
		s.Error <- errors.New("Lunarc is not running")
		s.Done <- false
	}
}

const aFilename = "access.log"

// LoggingServeMux logs HTTP requests
type LoggingServeMux struct {
	serveMux *http.ServeMux
	conf     Config
}

// NewLoggingServeMux allocates and returns a new LoggingServeMux
func NewLoggingServeMux(conf Config) *LoggingServeMux {
	serveMux := http.NewServeMux()
	return &LoggingServeMux{serveMux, conf}
}

// Handler sastisfy interface
func (mux *LoggingServeMux) Handler(r *http.Request) (h http.Handler, pattern string) {
	return mux.serveMux.Handler(r)
}

//ServeHTTP
func (mux *LoggingServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.serveMux.ServeHTTP(w, r)
}

//Handle register handler
func (mux *LoggingServeMux) Handle(pattern string, handler http.Handler) {

	var log = logrus.New()

	logFile, err := os.OpenFile(mux.conf.Log.File+aFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Out = os.Stderr
		log.Warningf("Can't open logfile: %v", err)
	} else {
		log.Out = logFile
	}
	mux.serveMux.Handle(pattern, Logging(handler, log))
}

// HandleFunc registers the handler function for the given pattern.
func (mux *LoggingServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.serveMux.Handle(pattern, http.HandlerFunc(handler))
}
