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
	"gopkg.in/mgo.v2"

	"github.com/DamienFontaine/lunarc/config"
)

//Context of a Server
type Context interface {
	GetCnf() config.Config
}

//DefaultContext DefaultContext
type DefaultContext struct {
	Cnf config.Config
}

//GetCnf returns config.Config
func (dc DefaultContext) GetCnf() config.Config {
	return dc.Cnf
}

//MongoContext add Mongo session to a Context
type MongoContext struct {
	Cnf     config.Config
	Session *mgo.Session
}

//GetCnf returns config.Config
func (mc MongoContext) GetCnf() config.Config {
	return mc.Cnf
}
