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
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

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
	Config    Config
	Error     chan error
	Done      chan bool
	quit      chan bool
	isStarted bool
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

	server = &Server{Config: conf, Done: make(chan bool, 1), Error: make(chan error, 1), Server: http.Server{Handler: NewLoggingServeMux(conf)}, quit: make(chan bool), isStarted: false}
	return
}

//Start the server.
func (s *Server) Start() (err error) {
	log.Infof("Lunarc is starting on port :%d", s.Config.Port)
	var l net.Listener
	go func() {
		l, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Config.Port))
		if err != nil {
			log.Errorf("Error: %v", err)
			s.Error <- err
			return
		}
		s.isStarted = true
		if len(s.Config.SSL.Certificate) > 0 && len(s.Config.SSL.Key) > 0 {
			err = s.ServeTLS(l, s.Config.SSL.Certificate, s.Config.SSL.Key)
			if err != nil && err != http.ErrServerClosed {
				log.Errorf("%v", err)
				l.Close()
				s.Error <- err
				s.quit <- true
			}
			close(s.quit)
		} else {
			err = s.Serve(l)
			if err != nil && err != http.ErrServerClosed {
				log.Errorf("%v", err)
				s.Error <- err
				s.quit <- true
			}
			close(s.quit)
		}
	}()

	<-s.quit

	if err = s.Shutdown(context.Background()); err != nil {
		log.Errorf("%v", err)
		s.Error <- err
	}

	<-s.quit

	l = nil
	log.Info("Lunarc terminated.")
	s.isStarted = false
	s.Done <- true
	return
}

//Stop the server.
func (s *Server) Stop() {
	if s.isStarted && s.quit != nil {
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
