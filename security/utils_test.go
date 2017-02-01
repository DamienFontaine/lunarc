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

package security

import (
	"bytes"
	"testing"
)

var result []byte
var salt []byte
var password []byte

func SetUp() {
	password = []byte("admin")
	salt = []byte{11, 155, 10, 202, 66, 120, 11, 199, 237, 220, 4, 25, 231, 174, 152, 247, 238, 34, 209, 179, 138, 116, 209, 132, 144, 114, 148, 175, 217, 157, 64, 236}
	result = []byte{174, 91, 45, 163, 0, 14, 75, 3, 235, 36, 79, 61, 159, 123, 165, 134, 117, 28, 193, 86, 30, 70, 110, 45, 83, 94, 248, 15, 51, 38, 120, 26}
}

func TestHashPassword(t *testing.T) {
	SetUp()
	hash, err := HashPassword(password, salt)
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}
	if !bytes.Equal(hash, result) {
		t.Fatalf("Non expected hash %v, need %v", hash, result)
	}
	if len(hash) != 32 {
		t.Fatalf("Incorrect hash size %v, need %v", len(hash), 32)
	}
}

func TestGenerateSalt(t *testing.T) {
	SetUp()
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}
	if len(salt) != 32 {
		t.Fatalf("Incorrect salt size %v, need %v", len(salt), 32)
	}
}

func TestCheckPasswordWithGoodPassword(t *testing.T) {
	SetUp()
	sGoodPassword := "admin"
	valid, err := CheckPassword([]byte(sGoodPassword), salt, result)
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}
	if !valid {
		t.Fatalf("Non expected behavior %v", sGoodPassword)
	}
}

func TestCheckPasswordWithBadPassword(t *testing.T) {
	SetUp()
	sBadPassword := "admi"
	valid, err := CheckPassword([]byte(sBadPassword), salt, result)
	if err != nil {
		t.Fatalf("Non expected error %v", err)
	}
	if valid {
		t.Fatalf("Non expected behavior %v", sBadPassword)
	}
}
