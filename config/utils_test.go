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

import "testing"

func TestConstructWithNormalByte(t *testing.T) {
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
      database: test`

	var serverEnvironment ServerEnvironment
	i, err := Get([]byte(data), "test", &serverEnvironment)
	server := i.(Server)
	if server.Port != 8888 {
		t.Fatalf("Non expected server port: %v != %v", 8888, server.Port)
	}

	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
}

func TestConstructWithBadByte(t *testing.T) {
	var data = `
   mongo::
    /port: 27017
    host: localhost
    "database": test
  `

	var mongoEnvironment MongoEnvironment
	_, err := Get([]byte(data), "test", &mongoEnvironment)

	if err == nil {
		t.Fatalf("Expected error: %v", err)
	}
}

func TestConstructWithNormalFile(t *testing.T) {
	var mongoEnvironment MongoEnvironment
	i, err := Get("config.yml", "test", &mongoEnvironment)
	mongo := i.(Mongo)
	if mongo.Port != 27017 {
		t.Fatalf("Non expected server port: %v != %v", 8888, mongo.Port)
	}

	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
}

func TestConstructWithNonExistentFile(t *testing.T) {
	var mongoEnvironment MongoEnvironment
	_, err := Get("no-config.yml", "test", &mongoEnvironment)

	if err == nil {
		t.Fatalf("Expected error!")
	}
}
