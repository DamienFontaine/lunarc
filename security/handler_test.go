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
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DamienFontaine/lunarc/web"
	jwt "github.com/dgrijalva/jwt-go"
)

func TestTokenHandlerNormal(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)

	cnf := new(web.Config)

	next := web.SingleFile("robot.txt")

	w := httptest.NewRecorder()
	TokenHandler(next, *cnf).ServeHTTP(w, request)

	if w.Code != http.StatusOK {
		t.Fatalf("Non expected code: %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "test") {
		t.Fatalf("Non expected Body")
	}
}

func TestTokenHandlerWithGoodToken(t *testing.T) {
	cnf := new(web.Config)

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = "test"
	claims["email"] = "test@test.com"
	claims["id"] = "id"
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	tokenString, err := token.SignedString([]byte(cnf.Jwt.Key))
	if err != nil {
		log.Fatal("Fatal", err)
	}

	request, _ := http.NewRequest("POST", "/robot.txt", nil)
	request.Header.Set("Authorization", "bearer "+tokenString)

	next := web.SingleFile("robot.txt")

	w := httptest.NewRecorder()
	TokenHandler(next, *cnf).ServeHTTP(w, request)

	if w.Code != http.StatusOK {
		t.Fatalf("Non expected code: %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "test") {
		t.Fatalf("Non expected Body")
	}
}

func TestTokenHandlerWithBadToken(t *testing.T) {
	cnf := new(web.Config)

	token := jwt.New(jwt.SigningMethodRS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = "test"
	claims["email"] = "test@test.com"
	claims["id"] = "id"
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	tokenString, _ := token.SignedString([]byte(cnf.Jwt.Key))

	request, _ := http.NewRequest("POST", "/robot.txt", nil)
	request.Header.Set("Authorization", "bearer "+tokenString)

	next := web.SingleFile("robot.txt")

	w := httptest.NewRecorder()
	TokenHandler(next, *cnf).ServeHTTP(w, request)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Non expected code: %v", w.Code)
	}
}

func TestTokenHandler401Error(t *testing.T) {
	request, _ := http.NewRequest("POST", "/robot.txt", nil)

	cnf := new(web.Config)

	next := web.SingleFile("robot.txt")

	w := httptest.NewRecorder()
	TokenHandler(next, *cnf).ServeHTTP(w, request)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Non expected code: %v", w.Code)
	}
}

func TestOAuth2WithoutToken(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	cnf := new(web.Config)
	next := web.SingleFile("robot.txt")
	w := httptest.NewRecorder()
	Oauth2(next, *cnf).ServeHTTP(w, request)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Non expected code: %v", w.Code)
	}
}
func TestOAuth2WithGoodToken(t *testing.T) {
	cnf := new(web.Config)
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	tokenString, _ := token.SignedString([]byte(cnf.Jwt.Key))
	request, _ := http.NewRequest("POST", "/robot.txt", nil)
	request.Header.Set("Authorization", "bearer "+tokenString)
	next := web.SingleFile("robot.txt")
	w := httptest.NewRecorder()
	Oauth2(next, *cnf).ServeHTTP(w, request)
	if w.Code != http.StatusOK {
		t.Fatalf("Non expected code: %v", w.Code)
	}
}
