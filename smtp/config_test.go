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

package smtp

import (
	"strings"
	"testing"

	"github.com/DamienFontaine/lunarc/config"

	"gopkg.in/yaml.v2"
)

func TestGetSMTPEnvironmentNormal(t *testing.T) {
	var data = `
  development:
    smtp:
      port: 464
      host: smtp.test.com
      auth:
        user: joh
        password: doe
  production:
    smtp:
      port: 465
      host: smtp.test.com
      auth:
        user: john
        password: doe
  `

	var smtpEnvironment SMTPEnvironment
	err := yaml.Unmarshal([]byte(data), &smtpEnvironment)
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}
	if len(smtpEnvironment.Env) != 2 {
		t.Fatalf("Must return 3 environments but %v", len(smtpEnvironment.Env))
	}
	var production Config
	res := smtpEnvironment.GetEnvironment("production")
	if res == nil {
		t.Fatalf("Must return a Server for production environment")
	}
	production = res.(Config)
	if production.Port != 465 {
		t.Fatalf("Must return a server port")
	}
	if strings.Compare(production.Auth.User, "john") != 0 {
		t.Fatalf("Must return a user")
	}
}

func TestSMTPEnvironment(t *testing.T) {
	var data = `
  development:
    smtp:
      port: 8888
      jwt:
        key: LunarcSecretKey
  `
	var smtpEnvironment SMTPEnvironment
	_ = yaml.Unmarshal([]byte(data), &smtpEnvironment)
	var i interface{} = &smtpEnvironment
	_, ok := i.(config.Environment)

	if !ok {
		t.Fatalf("SMTPServerEnvironment must implement Environment")
	}
}

func TestGetSMTPServerNormal(t *testing.T) {
	var data = `
  development:
    smtp:
      port: 465
      jwt:
        key: LunarcSecretKey
  `
	smtp, err := GetSMTP([]byte(data), "development")
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}
	if smtp.Port != 465 {
		t.Fatalf("Must return a Server with Port 465 not %v", smtp.Port)
	}
}
