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
	"errors"
	"net/smtp"
	"reflect"
	"testing"

	"github.com/DamienFontaine/lunarc/config"
	"github.com/DamienFontaine/lunarc/utils"
)

var mailService MailService

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

func TestNewMailServiceNormal(t *testing.T) {
	var data = `
  staging:
    smtp:
      port: 25
      host: smtp.doe.com
      auth:
        user: john@doe.com
        password: doe
  `
	server, _ := config.GetSMTP([]byte(data), "staging")
	mailService := NewMailService(server)
	f1 := reflect.ValueOf(smtp.SendMail)
	f2 := reflect.ValueOf(mailService.send)
	if f1.Pointer() != f2.Pointer() {
		t.Fatalf("MailService without SSL must use smtp.SendMail")
	}
}

func TestNewMailServiceWithSSL(t *testing.T) {
	var data = `
  staging:
    smtp:
      port: 465
      host: smtp.doe.com
      ssl: true
      auth:
        user: john@doe.com
        password: doe
  `
	server, _ := config.GetSMTP([]byte(data), "staging")
	mailService := NewMailService(server)
	f1 := reflect.ValueOf(utils.SendMailSSL)
	f2 := reflect.ValueOf(mailService.send)
	if f1.Pointer() != f2.Pointer() {
		t.Fatalf("MailService without SSL must use smtp.SendMail")
	}
}

func TestMailService(t *testing.T) {
	mailService = MailService{}

	var i interface{} = &mailService
	_, ok := i.(IMailService)

	if !ok {
		t.Fatalf("MailService must implement IMailService")
	}
}

func TestSendNormal(t *testing.T) {
	var data = `
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

func TestSendError(t *testing.T) {
	var data = `
  staging:
    smtp:
      port: 465
      host: smtp.doe.com
      ssl: true
      auth:
        user: john@doe.com
        password: doe
  `
	err := errors.New("Error")
	f, _ := mockSend(err)
	smtp, _ := config.GetSMTP([]byte(data), "staging")
	mailService = MailService{SMTP: smtp, send: f}

	err = mailService.Send("message", "test", "john@doe.com", "jane@doe.com")
	if err == nil {
		t.Fatalf("Must return an error")
	}
}
