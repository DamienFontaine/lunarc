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
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	jose "gopkg.in/square/go-jose.v2"

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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

//RandStringBytesMaskImprSrc Generate a random string
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

//EncodeOAuth2Code generate an OAuth2 code
func EncodeOAuth2Code(clientID, redirectURI, userID, sharedKey string) (code string, err error) {
	rand := RandStringBytesMaskImprSrc(20)
	exp := time.Now().Add(time.Minute * 10).String()
	response := NewResponse(clientID, redirectURI, userID, exp, rand)
	jresponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	j64response := base64.StdEncoding.EncodeToString(jresponse)
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS512, Key: []byte(sharedKey)}, nil)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	object, err := signer.Sign([]byte(j64response))
	if err != nil {
		log.Printf("Error: %v", err)
	}
	code, err = object.CompactSerialize()
	return
}

//DecodeOAuth2Code inverse of EncodeOAuth2Code
func DecodeOAuth2Code(code, sharedKey string) (response Response, err error) {
	object, err := jose.ParseSigned(code)
	if err != nil {
		return
	}
	output, err := object.Verify([]byte(sharedKey))
	if err != nil {
		return
	}
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(output)))
	l, err := base64.StdEncoding.Decode(base64Text, output)
	if err != nil {
		return
	}
	response = Response{}
	err = json.Unmarshal(base64Text[:l], &response)
	return
}
