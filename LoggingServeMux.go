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
	"os"

	"github.com/DamienFontaine/lunarc/config"

	"github.com/Sirupsen/logrus"
)

const aFilename = "access.log"

// LoggingServeMux logs HTTP requests
type LoggingServeMux struct {
	serveMux *http.ServeMux
	conf     config.Server
}

// NewLoggingServeMux allocates and returns a new LoggingServeMux
func NewLoggingServeMux(conf config.Server) *LoggingServeMux {
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
