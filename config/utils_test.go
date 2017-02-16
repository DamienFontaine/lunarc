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
)

//TestConfig
type TestConfig struct {
	Port int
	Log  struct {
		File  string
		Level string
	}
}

//TestEnvironment configurations
type TestEnvironment struct {
	Env map[string]TestConfig
}

//UnmarshalYAML implements Unmarshaler. Avoid use of env in the YAML file.
func (se *TestEnvironment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Env map[string]struct {
			TestConfig TestConfig
		}
	}
	if err := unmarshal(&aux.Env); err != nil {
		return err
	}
	se.Env = make(map[string]TestConfig)
	for env, conf := range aux.Env {
		se.Env[env] = conf.TestConfig
	}
	return nil
}

//GetEnvironment returns a Server configuration for the specified environment in parameter
func (se *TestEnvironment) GetEnvironment(environment string) interface{} {
	for env, conf := range se.Env {
		if strings.Compare(environment, env) == 0 {
			return conf
		}
	}
	return nil
}

func TestConstructWithNormalByte(t *testing.T) {
	var data = `
  development:
    testconfig:
      port: 8888
    source:
      port: 27017
      host: localhost
      database: test
  test:
    testconfig:
      port: 8888
    source:
      port: 27017
      host: mongo
      database: test
  production:
    testconfig:
      port: 8888
    source:
      port: 27017
      host: mongo
      database: test`

	var testEnvironment TestEnvironment
	i, err := Get([]byte(data), "test", &testEnvironment)
	test := i.(TestConfig)
	if test.Port != 8888 {
		t.Fatalf("Non expected server port: %v != %v", 8888, test.Port)
	}

	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
}

func TestConstructWithBadByte(t *testing.T) {
	var data = `
   testconfig::
    /port: 27017
    host: localhost
    "database": test
  `

	var testEnvironment TestEnvironment
	_, err := Get([]byte(data), "test", &testEnvironment)

	if err == nil {
		t.Fatalf("Expected error: %v", err)
	}
}

func TestConstructWithNormalFile(t *testing.T) {
	var testEnvironment TestEnvironment
	i, err := Get("config.yml", "test", &testEnvironment)
	test := i.(TestConfig)
	if test.Port != 8888 {
		t.Fatalf("Non expected server port: %v != %v", 8888, test.Port)
	}

	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
}

func TestConstructWithNonExistentFile(t *testing.T) {
	var testEnvironment TestEnvironment
	_, err := Get("no-config.yml", "test", &testEnvironment)

	if err == nil {
		t.Fatalf("Expected error!")
	}
}
