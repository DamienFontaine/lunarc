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

//SMTP configuration
type SMTP struct {
	Port int
	Host string
	SSL  bool
	Auth struct {
		User     string
		Password string
	}
}

//SMTPEnvironment configurations
type SMTPEnvironment struct {
	Env map[string]SMTP
}

//UnmarshalYAML implements Unmarshaler. Avoid use of env in the YAML file.
func (se *SMTPEnvironment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Env map[string]struct {
			SMTP SMTP
		}
	}
	if err := unmarshal(&aux.Env); err != nil {
		return err
	}
	se.Env = make(map[string]SMTP)
	for env, conf := range aux.Env {
		se.Env[env] = conf.SMTP
	}
	return nil
}

//GetEnvironment returns a SMTP Server configuration for the specified environment in parameter
func (se *SMTPEnvironment) GetEnvironment(environment string) interface{} {
	for env, conf := range se.Env {
		if strings.Compare(environment, env) == 0 {
			return conf
		}
	}
	return nil
}

//GetSMTP returns a SMTP Server configurations
func GetSMTP(source interface{}, environment string) (smtp SMTP, err error) {
	var smtpEnvironment SMTPEnvironment
	i, err := Get(source, environment, &smtpEnvironment)
	smtp = i.(SMTP)
	return
}
