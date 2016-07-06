// +build integration

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

package services

import (
	"net/smtp"
	"testing"

	"github.com/DamienFontaine/lunarc/config"
)

var mailService MailService

func TestMailService(t *testing.T) {
	mailService = MailService{}

	var i interface{} = &mailService
	_, ok := i.(IMailService)

	if !ok {
		t.Fatalf("MailService must implement IMailService")
	}
}

type EmailRecorder struct {
	addr string
	auth smtp.Auth
	from string
	to   []string
	msg  []byte
}

func mockSend(errToReturn error) (func(string, smtp.Auth, string, []string, []byte) error, *EmailRecorder) {
	r := new(EmailRecorder)
	return func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		*r = EmailRecorder{addr, a, from, to, msg}
		return errToReturn
	}, r
}

func TestSendNormal(t *testing.T) {
	var data = `
  development:
    smtp:
      port: 464
      host: smtp.test.com
      auth:
        user: john@doe.com
        password: doe
  staging:
    smtp:
      port: 465
      host: smtp.doe.com
      ssl: true
      auth:
        user: john@doe.com
        password: doe
  `
	body := "From: john@doe.com\r\n" +
		"To: jane@doe.com\r\n" +
		"Subject: test\r\n" +
		"\r\n" +
		"message\r\n"
	f, r := mockSend(nil)
	smtp, _ := config.GetSMTP([]byte(data), "staging")
	mailService = MailService{SMTP: smtp, send: f}

	err := mailService.Send("message", "test", "john@doe.com", "jane@doe.com")
	if err != nil {
		t.Fatalf("Mustn't return an error")
	}
	if string(r.msg) != body {
		t.Errorf("wrong message body.\n\nexpected: %s\n got: %s", body, r.msg)
	}
}
