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
	"strings"

	"github.com/DamienFontaine/lunarc/config"
)

//Config of a web server
type Config struct {
	Port  int
	URL   string
	Admin string
	Log   struct {
		File  string
		Level string
	}
	SSL struct {
		Key         string
		Certificate string
	}
	Jwt struct {
		Key string
	}
}

//ServerEnvironment configurations
type ServerEnvironment struct {
	Env map[string]Config
}

//UnmarshalYAML implements Unmarshaler. Avoid use of env in the YAML file.
func (se *ServerEnvironment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Env map[string]struct {
			Config Config `yaml:"server"`
		}
	}
	if err := unmarshal(&aux.Env); err != nil {
		return err
	}
	se.Env = make(map[string]Config)
	for env, conf := range aux.Env {
		se.Env[env] = conf.Config
	}
	return nil
}

//GetEnvironment returns a Server configuration for the specified environment in parameter
func (se *ServerEnvironment) GetEnvironment(environment string) interface{} {
	for env, conf := range se.Env {
		if strings.Compare(environment, env) == 0 {
			return conf
		}
	}
	return nil
}

//GetConfig returns a Server configurations
func GetConfig(source interface{}, environment string) (server Config, err error) {
	var serverEnvironment ServerEnvironment
	i, err := config.Get(source, environment, &serverEnvironment)
	server = i.(Config)
	return
}
