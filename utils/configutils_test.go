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

	var util ConfigUtil
	config, err := util.Construct([]byte(data), "test")
	if config.Server.Port != 8888 {
		t.Fatalf("Non expected server port: %v != %v", 8888, config.Server.Port)
	}

	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
}

func TestConstructWithEmptyMongo(t *testing.T) {
	var data = `test:
  server:
    port: 8888
    jwt:
      key: LunarcSecretKey`

	var util ConfigUtil
	config, err := util.Construct([]byte(data), "test")

	if config.Server.Port != 8888 {
		t.Fatalf("Non expected server port: %v != %v", 8888, config.Server.Port)
	}

	if config.Mongo.Port != 0 {
		t.Fatalf("Non expected server port: %v != %v", 0, config.Mongo.Port)
	}

	if err == nil {
		t.Fatalf("Expected error: %v", err)
	}
}

func TestConstructWithEmptyServer(t *testing.T) {
	var data = `
  mongo:
    port: 27017
    host: localhost
    database: test
  `

	var util ConfigUtil
	config, err := util.Construct([]byte(data), "test")

	if config.Server.Port != 0 {
		t.Fatalf("Non expected server port: %v != %v", 0, config.Server.Port)
	}

	if err == nil {
		t.Fatalf("Expected error: %v", err)
	}
}

func TestConstructWithBadByte(t *testing.T) {
	var data = `
   mongo::
    /port: 27017
    host: localhost
    "database": test
  `

	var util ConfigUtil
	_, err := util.Construct([]byte(data), "test")

	if err == nil {
		t.Fatalf("Expected error: %v", err)
	}
}

func TestConstructWithNormalFile(t *testing.T) {
	var util ConfigUtil
	config, err := util.Construct("config.yml", "test")

	if config.Server.Port != 8888 {
		t.Fatalf("Non expected server port: %v != %v", 8888, config.Server.Port)
	}

	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
}

func TestConstructWithNonExistentFile(t *testing.T) {
	var util ConfigUtil
	_, err := util.Construct("no-config.yml", "test")

	if err == nil {
		t.Fatalf("Expected error!")
	}
}
