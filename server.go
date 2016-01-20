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
	"net/http"

	"github.com/DamienFontaine/lunarc/config"
	"github.com/DamienFontaine/lunarc/datasource"
	"github.com/DamienFontaine/lunarc/utils"
)

//Server is an http.ServeMux with a Context.
type Server interface {
	Initialize(string, string) error
	Start() error
	Stop()
	GetContext() Context
	GetMux() *http.ServeMux
}

//WebServer is a Server with a specialize Context.
type WebServer struct {
	Mux     *http.ServeMux
	Context MongoContext
}

//Initialize the server.
func (ws *WebServer) Initialize(filename string, environment string) (err error) {
	var cnf config.Config

	if ws.Mux == nil {
		ws.Mux = http.NewServeMux()
	}

	var configUtil = new(utils.ConfigUtil)
	cnf, err = configUtil.Construct(filename, environment)

	ws.Context = MongoContext{cnf, nil}

	return
}

//Start the server.
func (ws *WebServer) Start() (err error) {
	log.Println("Lunarc is starting...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", ws.Context.Cnf.Server.Port), ws.Mux)
	if err != nil {
		log.Printf("Fatal: %v", err)
	}
	return err
}

//Stop the server.
func (ws *WebServer) Stop() {
	//TODO
}

//SetDatasource allow use of datasource. Here MongoDB.
func (ws *WebServer) SetDatasource() {
	ws.Context.Session = datasource.GetSession(ws.Context.Cnf.Mongo.Port, ws.Context.Cnf.Mongo.Host)
}

//GetContext returns the Context of the server.
func (ws *WebServer) GetContext() Context {
	return MongoContext(ws.Context)
}

//GetMux return the http.ServeMux of the server.
func (ws *WebServer) GetMux() *http.ServeMux {
	return ws.Mux
}
