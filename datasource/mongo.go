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

package datasource

import "gopkg.in/mgo.v2"

//GetSession retourne une session MongoDB.
func GetSession(port int, host string) *mgo.Session {
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}

//Mongo is a datasource.
type Mongo struct {
	Session  *mgo.Session
	Database *mgo.Database
}

//Copy retourne une cope de la session
func (m *Mongo) Copy() *Mongo {
	copy := m.Session.Copy()
	return &Mongo{Session: copy, Database: m.Database}
}

//Close a Mongo session
func (m *Mongo) Close() {
	m.Session.Close()
}
