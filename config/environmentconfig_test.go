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
	"testing"

	"gopkg.in/yaml.v2"
)

func TestGetEnvironmentNormal(t *testing.T) {
	var data = `
  development:
    server:
      port: 8888
      jwt:
        key: LunarcSecretKey
    mongo:
      port: 27017
      host: localhost
      database: test
  test:
    server:
      port: 8888
      jwt:
        key: LunarcSecretKey
    mongo:
      port: 27017
      host: mongo
      database: test
  production:
    server:
      port: 8888
      jwt:
        key: LunarcSecretKey
    mongo:
      port: 27017
      host: mongo
      database: test
  `
	var environmentConfig EnvironmentConfig
	err := yaml.Unmarshal([]byte(data), &environmentConfig)
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}
	var production Config
	if environmentConfig.GetEnvironment(&production, "production"); &production == nil {
		t.Fatalf("Must return a config.Config for production environment")
	}

	if production.Mongo.Port != 27017 {
		t.Fatalf("Must return a mongo port")
	}

	var dev Config
	if environmentConfig.GetEnvironment(&dev, "dev"); &dev == nil {
		t.Fatalf("Musn't return a config.Config for dev")
	}
}

func TestUnmarshalYAMLNormal(t *testing.T) {
	var data = `
  development:
    server:
      port: 8888
      jwt:
        key: LunarcSecretKey
    mongo:
      port: 27017
      host: localhost
      database: test
  test:
    server:
      port: 8888
      jwt:
        key: LunarcSecretKey
    mongo:
      port: 27017
      host: mongo
      database: test
  production:
    server:
      port: 8888
      jwt:
        key: LunarcSecretKey
    mongo:
      port: 27017
      host: mongo
      database: test
  `
	var environmentConfig EnvironmentConfig
	err := yaml.Unmarshal([]byte(data), &environmentConfig)
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}

	if len(environmentConfig.Env) != 3 {
		t.Fatalf("Not enought configuration: %v in place of 2", len(environmentConfig.Env))
	}

	var dev Config
	dev = environmentConfig.Env["development"]
	if strings.Compare(dev.Mongo.Host, "localhost") != 0 {
		t.Fatalf("Dev Mongo Host: %v in place of localhost", dev.Mongo.Host)
	}

	var test Config
	test = environmentConfig.Env["test"]
	if strings.Compare(test.Mongo.Host, "mongo") != 0 {
		t.Fatalf("Test Mongo Host: %v in place of mongo", test.Mongo.Host)
	}
}
