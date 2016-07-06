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
	"fmt"
	"log"
	"net/smtp"

	"github.com/DamienFontaine/lunarc/config"
	"github.com/DamienFontaine/lunarc/utils"
)

//IMailService interface
type IMailService interface {
	Send(message string, subject string, from string, to string) error
}

//MailService send email
type MailService struct {
	SMTP config.SMTP
	send func(string, smtp.Auth, string, []string, []byte) error
}

//NewMailService retourne un MailService
func NewMailService(server config.SMTP) *MailService {
	f := smtp.SendMail
	if server.SSL {
		f = utils.SendMailSSL
	}
	return &MailService{SMTP: server, send: f}
}

//Send envoie un email
func (m *MailService) Send(message string, subject string, from string, to string) (err error) {
	auth := smtp.PlainAuth("", m.SMTP.Auth.User, m.SMTP.Auth.Password, m.SMTP.Host)

	t := []string{to}
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")

	err = m.send(fmt.Sprintf("%s:%d", m.SMTP.Host, m.SMTP.Port), auth, from, t, msg)
	if err != nil {
		log.Fatal(err)
	}

	return
}
