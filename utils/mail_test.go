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

package utils

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"testing"
)

var sendMailServer = `220 hello world
502 EH?
250 smtp.doe.com at your service
250 Sender ok
250 Receiver ok
354 Go ahead
250 Data ok
221 Goodbye
`

var sendMailClient = `EHLO localhost
HELO localhost
MAIL FROM:<john@doe.com>
RCPT TO:<jane@doe.com>
DATA
From: john@doe.com
To: jane@doe.com
Subject: SendMail test

SendMailSSL is working for me.
.
QUIT
`

func TestSendMailSSL(t *testing.T) {
	server := strings.Join(strings.Split(sendMailServer, "\n"), "\r\n")
	client := strings.Join(strings.Split(sendMailClient, "\n"), "\r\n")
	var cmdbuf bytes.Buffer
	bcmdbuf := bufio.NewWriter(&cmdbuf)
	cer, err := tls.LoadX509KeyPair("../testdata/ssl/smtp.crt", "../testdata/ssl/smtp.key")
	if err != nil {
		log.Println(err)
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	l, err := tls.Listen("tcp", "127.0.0.1:0", config)
	if err != nil {
		t.Fatalf("Unable to to create listener: %v", err)
	}
	defer l.Close()

	// prevent data race on bcmdbuf
	var done = make(chan struct{})
	go func(data []string) {

		defer close(done)
		var conn net.Conn
		conn, err = l.Accept()
		if err != nil {
			t.Errorf("Accept error: %v", err)
			return
		}
		defer conn.Close()

		tc := textproto.NewConn(conn)
		for i := 0; i < len(data) && data[i] != ""; i++ {
			tc.PrintfLine(data[i])
			for len(data[i]) >= 4 && data[i][3] == '-' {
				i++
				tc.PrintfLine(data[i])
			}
			if data[i] == "221 Goodbye" {
				return
			}
			read := false
			for !read || data[i] == "354 Go ahead" {
				var msg string
				msg, err = tc.ReadLine()
				bcmdbuf.Write([]byte(msg + "\r\n"))
				read = true
				if err != nil {
					t.Errorf("Read error: %v", err)
					return
				}
				if data[i] == "354 Go ahead" && msg == "." {
					break
				}
			}
		}
	}(strings.Split(server, "\r\n"))

	user := "john@doe.com"
	password := "doe"
	host := "smtp.doe.com"
	from := "john@doe.com"
	to := []string{"jane@doe.com"}
	msg := []byte("From: john@doe.com\r\n" +
		"To: jane@doe.com\r\n" +
		"Subject: SendMail test\r\n" +
		"\r\n" +
		"SendMailSSL is working for me.")

	a := smtp.PlainAuth("", user, password, host)
	err = SendMailSSL(l.Addr().String(), a, from, to, msg)
	if err != nil {
		t.Fatalf("Mustn't return error %s", err)
	}
	<-done
	bcmdbuf.Flush()
	actualcmds := cmdbuf.String()
	if client != actualcmds {
		t.Errorf("Got:\n%s\nExpected:\n%s", actualcmds, client)
	}
}
