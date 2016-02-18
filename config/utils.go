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
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Get a config
func Get(source interface{}, environment string, configEnv Environment) (conf interface{}, err error) {
	if filename, ok := source.(string); ok {
		source, err = ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Fatal: %v", err)
			return
		}
	}

	err = yaml.Unmarshal(source.([]byte), configEnv)
	if err != nil {
		log.Printf("Fatal: bad config : %v", err)
		return
	}

	conf = configEnv.GetEnvironment(environment)
	if conf == nil {
		err = errors.New("No configuration")
		return
	}
	return
}
