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
	"testing"
)

type EmailRecorder struct {
	from string
	to   []string
	msg  []byte
}

type SMTPMock struct {
	r   *EmailRecorder
	err error
}

func (s *SMTPMock) SendMail(from string, to []string, msg []byte) error {
	s.r = &EmailRecorder{from, to, msg}
	return s.err
}

func TestMailService(t *testing.T) {
	mailService := MailService{}

	var i interface{} = &mailService
	_, ok := i.(IMailService)

	if !ok {
		t.Fatalf("MailService must implement IMailService")
	}
}

func TestSendNormal(t *testing.T) {
	body := "From: john@doe.com\r\n" +
		"To: jane@doe.com\r\n" +
		"Subject: test\r\n" +
		"\r\n" +
		"message\r\n"
	s := SMTPMock{err: nil}
	mailService := NewMailService(&s)

	err := mailService.Send("message", "test", "john@doe.com", "jane@doe.com")
	if err != nil {
		t.Fatalf("Mustn't return an error")
	}
	if string(s.r.msg) != body {
		t.Errorf("wrong message body.\n\nexpected: %s\n got: %s", body, s.r.msg)
	}
}

func TestSendError(t *testing.T) {
	err := errors.New("Error")
	s := SMTPMock{err: err}
	mailService := NewMailService(&s)

	err = mailService.Send("message", "test", "john@doe.com", "jane@doe.com")
	if err == nil {
		t.Fatalf("Must return an error")
	}
}
