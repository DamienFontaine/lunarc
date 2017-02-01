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
	"testing"

	"github.com/DamienFontaine/lunarc/config"

	"gopkg.in/yaml.v2"
)

func TestGetMongoEnvironmentNormal(t *testing.T) {
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
      credential:
        username: lunarc
        password: lunarc
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

	var mongoEnvironment MongoEnvironment
	err := yaml.Unmarshal([]byte(data), &mongoEnvironment)
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}
	if len(mongoEnvironment.Env) != 3 {
		t.Fatalf("Must return 3 environments but %v", len(mongoEnvironment.Env))
	}
	var production Config
	res := mongoEnvironment.GetEnvironment("production")
	if res == nil {
		t.Fatalf("Must return a Mongo for production environment")
	}
	production = res.(Config)
	if production.Port != 27017 {
		t.Fatalf("Must return a mongo port")
	}
}

func TestMongoEnvironment(t *testing.T) {
	var data = `
  development:
    mongo:
      port: 27017
      host: mongo
      database: test
  `
	var mongoEnvironment MongoEnvironment
	_ = yaml.Unmarshal([]byte(data), &mongoEnvironment)
	var i interface{} = &mongoEnvironment
	_, ok := i.(config.Environment)

	if !ok {
		t.Fatalf("MongoEnvironment must implement Environment")
	}
}

func TestGetMongoNormal(t *testing.T) {
	var data = `
  development:
    mongo:
      port: 27017
      host: mongo
      database: test
      credential:
        username: lunarc
        password: lunarc
  `
	mongo, err := GetMongo([]byte(data), "development")
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}
	if mongo.Port != 27017 {
		t.Fatalf("Must return a Mongo with Port 27017 not %v", mongo.Port)
	}
	if strings.Compare(mongo.Host, "mongo") != 0 {
		t.Fatalf("Must return a Mongo with Host mongo not %v", mongo.Host)
	}
	if strings.Compare(mongo.Database, "test") != 0 {
		t.Fatalf("Must return a Mongo with Database mongo not %v", mongo.Database)
	}
	if strings.Compare(mongo.Credential.Username, "lunarc") != 0 {
		t.Fatalf("Must return a Mongo with username credential lunarc not %v", mongo.Credential.Username)
	}
	if strings.Compare(mongo.Credential.Password, "lunarc") != 0 {
		t.Fatalf("Must return a Mongo with password credential lunarc not %v", mongo.Credential.Password)
	}
}
