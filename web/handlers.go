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
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

//Logging logs http requests
func Logging(next http.Handler, log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := StatusResponseWriter{w, 0, 0}
		start := time.Now()
		next.ServeHTTP(&srw, r)
		end := time.Now()
		latency := end.Sub(start)

		log.WithField("client", r.RemoteAddr).WithField("latency", latency).WithField("length", srw.Length()).WithField("code", srw.Status()).Printf("%s %s %s", r.Method, r.URL, r.Proto)
	})
}

//SingleFile returns a handler
func SingleFile(filename string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

// StatusResponseWriter returns status code
type StatusResponseWriter struct {
	http.ResponseWriter
	status int
	length int
}

// Status return status code
func (w *StatusResponseWriter) Status() int {
	return w.status
}

// Length return response size
func (w *StatusResponseWriter) Length() int {
	return w.length
}

// Header Satisfy the http.ResponseWriter interface
func (w *StatusResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write Satisfy the http.ResponseWriter interface
func (w *StatusResponseWriter) Write(data []byte) (int, error) {
	w.length = len(data)
	return w.ResponseWriter.Write(data)
}

//Hijack Satisfy the http.ResponseWriter interface
func (w *StatusResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("Not a Hijacker")
}

// WriteHeader writes status code
func (w *StatusResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
