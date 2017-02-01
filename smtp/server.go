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
	"fmt"
	"net/smtp"
)

//MailSender interface
type MailSender interface {
	SendMail(from string, to []string, msg []byte) error
}

//SMTP SMTP server
type SMTP struct {
	addr string
	auth smtp.Auth
	send func(string, smtp.Auth, string, []string, []byte) error
}

//NewSMTP create new SMTP
func NewSMTP(filename string, environment string) (s *SMTP, err error) {
	conf, err := GetSMTP(filename, environment)
	if err != nil {
		return
	}
	auth := smtp.PlainAuth("", conf.Auth.User, conf.Auth.Password, conf.Host)
	f := smtp.SendMail
	if conf.SSL {
		f = SendMailSSL
	}
	s = &SMTP{auth: auth, send: f, addr: fmt.Sprintf("%s:%d", conf.Host, conf.Port)}
	return
}

//SendMail send an email
func (s *SMTP) SendMail(from string, to []string, msg []byte) (err error) {
	err = s.send(s.addr, s.auth, from, to, msg)
	return
}
