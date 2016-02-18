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
	"testing"

	"gopkg.in/yaml.v2"
)

func TestGetServerEnvironmentNormal(t *testing.T) {
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

	var serverEnvironment ServerEnvironment
	err := yaml.Unmarshal([]byte(data), &serverEnvironment)
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}
	if len(serverEnvironment.Env) != 3 {
		t.Fatalf("Must return 3 environments but %v", len(serverEnvironment.Env))
	}
	var production Server
	res := serverEnvironment.GetEnvironment("production")
	if res == nil {
		t.Fatalf("Must return a Server for production environment")
	}
	production = res.(Server)
	if production.Port != 8888 {
		t.Fatalf("Must return a server port")
	}
}

func TestServerEnvironment(t *testing.T) {
	var data = `
  development:
    server:
      port: 8888
      jwt:
        key: LunarcSecretKey
  `
	var serverEnvironment ServerEnvironment
	_ = yaml.Unmarshal([]byte(data), &serverEnvironment)
	var i interface{} = &serverEnvironment
	_, ok := i.(Environment)

	if !ok {
		t.Fatalf("ServerEnvironment must implement Environment")
	}
}

func TestGetServerNormal(t *testing.T) {
	var data = `
  development:
    server:
      port: 8888
      jwt:
        key: LunarcSecretKey
  `
	server, err := GetServer([]byte(data), "development")
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}
	if server.Port != 8888 {
		t.Fatalf("Must return a Server with Port 8888 not %v", server.Port)
	}
}
