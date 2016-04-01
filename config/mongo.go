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

package config

import (
	"strings"

	"gopkg.in/mgo.v2"
)

//Mongo configuration
type Mongo struct {
	Port       int
	Host       string
	Database   string
	Credential *mgo.Credential
}

//MongoEnvironment configurations
type MongoEnvironment struct {
	Env map[string]Mongo
}

//UnmarshalYAML implements Unmarshaler. Avoid use of env in the YAML file.
func (m *MongoEnvironment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Env map[string]struct {
			Mongo
		}
	}
	if err := unmarshal(&aux.Env); err != nil {
		return err
	}
	m.Env = make(map[string]Mongo)
	for env, conf := range aux.Env {
		m.Env[env] = conf.Mongo
	}
	return nil
}

//GetEnvironment returns a Mongo configuration for the specified environment in parameter
func (m *MongoEnvironment) GetEnvironment(environment string) interface{} {
	for env, conf := range m.Env {
		if strings.Compare(environment, env) == 0 {
			return conf
		}
	}
	return nil
}

//GetMongo returns a Mongo configurations
func GetMongo(source interface{}, environment string) (mongo Mongo, err error) {
	var mongoEnvironment MongoEnvironment
	i, err := Get(source, environment, &mongoEnvironment)
	mongo = i.(Mongo)
	return
}
