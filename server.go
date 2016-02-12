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
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/DamienFontaine/lunarc/config"
	"github.com/DamienFontaine/lunarc/datasource"
	"github.com/DamienFontaine/lunarc/utils"
)

//Server is an http.ServeMux with a Context.
type Server interface {
	Start() error
	Stop()
	GetContext() Context
	GetHandler() http.Handler
}

//WebServer is a Server with a specialize Context.
type WebServer struct {
	Context   MongoContext
	server    http.Server
	quit      chan bool
	err       chan error
	interrupt chan os.Signal
}

//NewWebServer create a new instance of WebServer
func NewWebServer(filename string, environment string) (server *WebServer, err error) {
	var cnf config.Config
	var configUtil = new(utils.ConfigUtil)
	cnf, err = configUtil.Construct(filename, environment)
	context := MongoContext{cnf, nil}

	server = &WebServer{Context: context, server: http.Server{Handler: http.NewServeMux()}, quit: make(chan bool), err: make(chan error, 1)}
	return
}

//Start the server.
func (ws *WebServer) Start() (err error) {
	log.Println("Lunarc is starting...")
	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", ws.Context.Cnf.Server.Port))
		if err != nil {
			log.Printf("Error: %v", err)
			ws.err <- err
			return
		}
		go ws.handleInterrupt(l)
		err = ws.server.Serve(l)
		if err != nil {
			log.Printf("Error: %v", err)
			ws.interrupt <- syscall.SIGINT
			ws.err <- err
			return
		}
	}()

	for {
		select {
		case <-ws.quit:
			ws.interrupt <- syscall.SIGINT
			return
		default:
			//continue
		}
	}
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
	} else {
		log.Println("Lunarc is not running")
	}
}

//SetDatasource allow use of datasource. Here MongoDB.
func (ws *WebServer) SetDatasource() {
	ws.Context.Session = datasource.GetSession(ws.Context.Cnf.Mongo.Port, ws.Context.Cnf.Mongo.Host)
}

//GetContext returns the Context of the server.
func (ws *WebServer) GetContext() Context {
	return MongoContext(ws.Context)
}

//GetHandler return the http.ServeMux of the server.
func (ws *WebServer) GetHandler() http.Handler {
	return ws.server.Handler
}
