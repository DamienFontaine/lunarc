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

package utils

import (
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/DamienFontaine/lunarc/config"
)

var util ConfigUtil

//ConfigUtil helps to manipulate config.Config
type ConfigUtil struct{}

// Construct a config.Config
func (cu *ConfigUtil) Construct(source interface{}, environment string) (conf config.Config, err error) {
	var environmentConfig config.EnvironmentConfig

	if filename, ok := source.(string); ok {
		source, err = ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Fatal: %v", err)
			return
		}
	}

	err = yaml.Unmarshal(source.([]byte), &environmentConfig)
	if err != nil {
		log.Printf("Fatal: bad config : %v", err)
		return
	}

	environmentConfig.GetEnvironment(&conf, environment)
	if &conf == nil {
		err = errors.New("No configuration")
		return
	}

	if conf.Mongo.Port == 0 {
		log.Printf("Server Mongo misconfigured")
		err = errors.New("Server Mongo misconfigured")
		return
	}

	if conf.Server.Port == 0 {
		log.Printf("Server misconfigured")
		err = errors.New("Server misconfigured")
		return
	}
	return
}
