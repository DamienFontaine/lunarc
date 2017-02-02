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

package controllers

import "net/http"

//PingController ping.
type PingController struct {
}

//NewPingController create a new PingController
func NewPingController() *PingController {
	pingController := PingController{}
	return &pingController
}

//Ping respond pong
func (c *PingController) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
	return
}
