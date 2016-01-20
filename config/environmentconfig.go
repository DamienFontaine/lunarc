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

import "strings"

//EnvironmentConfig contains all configs
type EnvironmentConfig struct {
	Env map[string]Config
}

//UnmarshalYAML implements Unmarshaler. Avoid use of env in the YAML file.
func (ec *EnvironmentConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&ec.Env)
}

//GetEnvironment returns a config.Config for the specified environment in parameter
func (ec *EnvironmentConfig) GetEnvironment(config *Config, environment string) {
	for env, conf := range ec.Env {
		if strings.Compare(environment, env) == 0 {
			*config = conf
		}
	}
	config = nil
}
