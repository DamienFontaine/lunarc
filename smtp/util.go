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
	"crypto/tls"
	"log"
	"net"
	"net/smtp"
)

// SendMailSSL envoie un email par SSL
func SendMailSSL(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true}) //TODO: Not secure
	if err != nil {
		log.Println("Error Dialing", err)
		return err
	}
	h, _, _ := net.SplitHostPort(addr)
	c, err := smtp.NewClient(conn, h)
	if err != nil {
		log.Println("Error SMTP connection", err)
		return err
	}
	defer c.Close()

	if a != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(a); err != nil {
				log.Printf("Authentication error: %v", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		log.Printf("From error: %v", err)
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			log.Printf("Recipient error: %v", err)
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	w.Write(msg)
	w.Close()

	return c.Quit()
}
