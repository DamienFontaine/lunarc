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
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/DamienFontaine/lunarc/config"
)

//Server is an http.ServeMux with a Context.
type Server interface {
	Start() error
	Stop()
	GetHandler() http.Handler
	GetConfig() config.Server
}

//WebServer is a Server with a specialize Context.
type WebServer struct {
	Error     chan error
	Done      chan bool
	server    http.Server
	quit      chan bool
	interrupt chan os.Signal
	conf      config.Server
}

//NewWebServer create a new instance of WebServer
func NewWebServer(filename string, environment string) (server *WebServer, err error) {
	conf, err := config.GetServer(filename, environment)
	server = &WebServer{conf: conf, Done: make(chan bool, 1), Error: make(chan error, 1), server: http.Server{Handler: http.NewServeMux()}, quit: make(chan bool)}
	return
}

//Start the server.
func (ws *WebServer) Start() (err error) {
	log.Printf("Lunarc is starting on port :%d", ws.conf.Port)
	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", ws.conf.Port))
		if err != nil {
			log.Printf("Error: %v", err)
			ws.Error <- err
			return
		}
		go ws.handleInterrupt(l)
		err = ws.server.Serve(l)
		if err != nil {
			log.Printf("Error: %v", err)
			ws.interrupt <- syscall.SIGINT
			ws.Error <- err
			return
		}
	}()

	<-ws.quit
	ws.interrupt <- syscall.SIGINT
	return
}

func (ws *WebServer) handleInterrupt(listener net.Listener) {
	if ws.interrupt == nil {
		ws.interrupt = make(chan os.Signal, 1)
	}
	<-ws.interrupt
	listener.Close()
	ws.quit <- true
	ws.interrupt = nil
}

//Stop the server.
func (ws *WebServer) Stop() {
	if ws.interrupt != nil && ws.quit != nil {
		log.Println("Lunarc is stopping...")
		ws.quit <- true
		<-ws.quit
		log.Println("Lunarc stopped.")
		ws.Done <- true
	} else {
		log.Println("Lunarc is not running")
		ws.Error <- errors.New("Lunarc is not running")
	}
}

//GetHandler return the http.ServeMux of the server.
func (ws *WebServer) GetHandler() http.Handler {
	return ws.server.Handler
}

//GetConfig return the configuration of the server.
func (ws *WebServer) GetConfig() config.Server {
	return ws.conf
}
