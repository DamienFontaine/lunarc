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
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/scrypt"
)

//Constantes nécessaires à l'encryption du mot de passe
const (
	KEYLENGTH = 32
	N         = 16384
	R         = 8
	P         = 1
	S         = 32
)

//HashPassword hash un mot de passe
func HashPassword(password []byte, salt []byte) (hash []byte, err error) {
	hash, err = scrypt.Key(password, salt, N, R, P, KEYLENGTH)
	if err != nil {
		return nil, err
	}
	return
}

//GenerateSalt génère le salage
func GenerateSalt() (salt []byte, err error) {
	salt = make([]byte, S)
	_, err = rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return
}

//CheckPassword vérifie si le mot de passe est correct
func CheckPassword(password []byte, salt []byte, hpassword []byte) (bool, error) {
	hash, err := HashPassword(password, salt)
	if err != nil {
		return false, err
	}
	return bytes.Equal(hash, hpassword), nil
}
