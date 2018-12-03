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
	"strings"

	"github.com/DamienFontaine/lunarc/config"
)

//Config of Mongo
type Config struct {
	Port     int
	Host     string
	Database string
	Username string
	Password string
}

//Environment configurations
type Environment struct {
	Env map[string]Config
}

//UnmarshalYAML implements Unmarshaler. Avoid use of env in the YAML file.
func (m *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Env map[string]struct {
			Config `yaml:"mongo"`
		}
	}
	if err := unmarshal(&aux.Env); err != nil {
		return err
	}
	m.Env = make(map[string]Config)
	for env, conf := range aux.Env {
		m.Env[env] = conf.Config
	}
	return nil
}

//GetEnvironment returns a Mongo configuration for the specified environment in parameter
func (m *Environment) GetEnvironment(environment string) interface{} {
	for env, conf := range m.Env {
		if strings.Compare(environment, env) == 0 {
			return conf
		}
	}
	return nil
}

//GetMongo returns a Mongo configurations
func GetMongo(source interface{}, environment string) (mongo Config, err error) {
	var env Environment
	i, err := config.Get(source, environment, &env)
	mongo = i.(Config)
	return
}
