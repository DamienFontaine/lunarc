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

package mongo

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
)

//Mongo is a datasource.
type Mongo struct {
	Session  *mgo.Session
	Database *mgo.Database
}

//NewMongo creates a newinstance of Mongo
func NewMongo(filename string, environment string) (*Mongo, error) {
	cnf, err := GetMongo(filename, environment)
	if err != nil {
		return nil, err
	}
	session, err := mgo.Dial(fmt.Sprintf("%v:%d", cnf.Host, cnf.Port))
	if err != nil {
		log.Printf("Impossible de contacter %v sur le port %d", cnf.Host, cnf.Port)
		return nil, err
	}
	if cnf.Credential != nil {
		if len(cnf.Credential.Username) > 0 && len(cnf.Credential.Password) > 0 {
			err = session.Login(cnf.Credential)
			if err != nil {
				log.Println("Impossible de s'identifier à MongoDB")
				return nil, err
			}
		}
	}
	session.SetMode(mgo.Monotonic, true)
	return &Mongo{Session: session, Database: session.DB(cnf.Database)}, nil
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
