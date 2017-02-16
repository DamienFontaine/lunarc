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
250-smtp.doe.com at your service
250 Cool
250 Sender ok
250 Receiver ok
354 Go ahead
250 Data ok
221 Goodbye
`

var sendMailServerWithoutFrom = `220 hello world
502 EH?
250-smtp.doe.com at your service
250 Cool
451 Requested action aborted
250 Data ok
221 Goodbye
`

var sendMailServerWithoutRecipient = `220 hello world
502 EH?
250-smtp.doe.com at your service
250 Cool
250 Sender ok
250 Receiver ok
451 Requested action aborted
250 Data ok
221 Goodbye
`

var sendMailServerWithoutData = `220 hello world
502 EH?
250-smtp.doe.com at your service
250 Cool
250 Sender ok
451 Requested action aborted
250 Data ok
221 Goodbye
`

var sendMailServerWithAuth = `220 hello world
250-smtp.doe.com at your service
250 AUTH LOGIN PLAIN
221 Goodbye
`

var sendMailServerNotAvailable = `421 Not Available`

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

func FakeServer(data []string, bcmdbuf *bufio.Writer, done chan struct{}, l net.Listener, t *testing.T) {

	defer close(done)
	var conn net.Conn
	conn, err := l.Accept()
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
		if data[i] == "421 Not Available" {
			return
		}
		if data[i] == "535 Invalid credentials" {
			return
		}
		if data[i] == "451 Requested action aborted" {
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
}

func TestSendMailWithoutServer(t *testing.T) {
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
	l.Close()

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
	if err == nil {
		t.Fatal("Must return an erro")
	}
}

func TestSendMailWithNotAvailableError(t *testing.T) {
	server := strings.Join(strings.Split(sendMailServerNotAvailable, "\n"), "\r\n")
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
	go FakeServer(strings.Split(server, "\r\n"), bcmdbuf, done, l, t)

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
	if err == nil {
		t.Fatal("Mus return an error")
	}
}

func TestSendMailWithAuthenticationError(t *testing.T) {
	server := strings.Join(strings.Split(sendMailServerWithAuth, "\n"), "\r\n")
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
	go FakeServer(strings.Split(server, "\r\n"), bcmdbuf, done, l, t)

	user := "joh@doe.com"
	password := "does"
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
	if err == nil {
		t.Fatalf("Must return an error")
	}
}

func TestSendMailWithFromError(t *testing.T) {
	server := strings.Join(strings.Split(sendMailServerWithoutFrom, "\n"), "\r\n")
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
	go FakeServer(strings.Split(server, "\r\n"), bcmdbuf, done, l, t)

	user := "joh@doe.com"
	password := "does"
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
	if err == nil {
		t.Fatalf("Must return an error")
	}
}

func TestSendMailWithRecipientError(t *testing.T) {
	server := strings.Join(strings.Split(sendMailServerWithoutRecipient, "\n"), "\r\n")
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
	go FakeServer(strings.Split(server, "\r\n"), bcmdbuf, done, l, t)

	user := "joh@doe.com"
	password := "does"
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
	if err == nil {
		t.Fatalf("Must return an error")
	}
}

func TestSendMailWithDataError(t *testing.T) {
	server := strings.Join(strings.Split(sendMailServerWithoutData, "\n"), "\r\n")
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
	go FakeServer(strings.Split(server, "\r\n"), bcmdbuf, done, l, t)

	user := "joh@doe.com"
	password := "does"
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
	if err == nil {
		t.Fatalf("Must return an error")
	}
}

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
	go FakeServer(strings.Split(server, "\r\n"), bcmdbuf, done, l, t)

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
