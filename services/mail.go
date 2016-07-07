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

import "github.com/DamienFontaine/lunarc"

//IMailService interface
type IMailService interface {
	Send(message string, subject string, from string, to string) error
}

//MailService send email
type MailService struct {
	SMTP lunarc.MailSender
}

//NewMailService retourne un MailService
func NewMailService(server lunarc.MailSender) *MailService {
	return &MailService{SMTP: server}
}

//Send envoie un email
func (m *MailService) Send(message string, subject string, from string, to string) (err error) {
	t := []string{to}
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")

	err = m.SMTP.SendMail(from, t, msg)

	return
}
