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

package lunarc

import (
	"net/smtp"
	"reflect"
	"testing"

	"github.com/DamienFontaine/lunarc/utils"
)

func TestNewSMTPNormal(t *testing.T) {
	s, err := NewSMTP("config.yml", "test")
	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
	f1 := reflect.ValueOf(smtp.SendMail)
	f2 := reflect.ValueOf(s.send)
	if f1.Pointer() != f2.Pointer() {
		t.Fatalf("SMTP without SSL must use smtp.SendMail")
	}
}

func TestNewSMTPWithSSL(t *testing.T) {
	s, err := NewSMTP("config.yml", "ssl")
	if err != nil {
		t.Fatalf("Non expected error: %v", err)
	}
	f1 := reflect.ValueOf(utils.SendMailSSL)
	f2 := reflect.ValueOf(s.send)
	if f1.Pointer() != f2.Pointer() {
		t.Fatalf("SMTP with SSL must use utils.SendMailSSL")
	}
}
