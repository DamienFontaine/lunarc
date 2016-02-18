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
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/DamienFontaine/lunarc/config"
)

func TestSingleFileNormal(t *testing.T) {
	request, _ := http.NewRequest("GET", "robot.txt", nil)

	w := httptest.NewRecorder()
	SingleFile("robot.txt").ServeHTTP(w, request)

	if w.Code != http.StatusOK {
		t.Fatalf("Non expected code: %v", w.Code)
	}
}

func TestSingleFileNotFound(t *testing.T) {
	request, _ := http.NewRequest("GET", "robot.txt", nil)

	w := httptest.NewRecorder()
	SingleFile("robots.txt").ServeHTTP(w, request)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Non expected code: %v", w.Code)
	}
}

func TestAuthMiddleWareNormal(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)

	cnf := new(config.Server)

	next := SingleFile("robot.txt")

	w := httptest.NewRecorder()
	AuthMiddleWare(next, *cnf).ServeHTTP(w, request)

	if w.Code != http.StatusOK {
		t.Fatalf("Non expected code: %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "test") {
		t.Fatalf("Non expected Body")
	}
}

func TestAuthMiddleWareWithGoodToken(t *testing.T) {
	cnf := new(config.Server)

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["username"] = "test"
	token.Claims["email"] = "test@test.com"
	token.Claims["id"] = "id"
	token.Claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	tokenString, err := token.SignedString([]byte(cnf.Jwt.Key))
	if err != nil {
		log.Fatal("Fatal", err)
	}

	request, _ := http.NewRequest("POST", "/robot.txt", nil)
	request.Header.Set("Authorization", "bearer "+tokenString)

	next := SingleFile("robot.txt")

	w := httptest.NewRecorder()
	AuthMiddleWare(next, *cnf).ServeHTTP(w, request)

	if w.Code != http.StatusOK {
		t.Fatalf("Non expected code: %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "test") {
		t.Fatalf("Non expected Body")
	}
}

func TestAuthMiddleWareWithBadToken(t *testing.T) {
	cnf := new(config.Server)

	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims["username"] = "test"
	token.Claims["email"] = "test@test.com"
	token.Claims["id"] = "id"
	token.Claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	tokenString, _ := token.SignedString([]byte(cnf.Jwt.Key))

	request, _ := http.NewRequest("POST", "/robot.txt", nil)
	request.Header.Set("Authorization", "bearer "+tokenString)

	next := SingleFile("robot.txt")

	w := httptest.NewRecorder()
	AuthMiddleWare(next, *cnf).ServeHTTP(w, request)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Non expected code: %v", w.Code)
	}
}

func TestAuthMiddleWare401Error(t *testing.T) {
	request, _ := http.NewRequest("POST", "/robot.txt", nil)

	cnf := new(config.Server)

	next := SingleFile("robot.txt")

	w := httptest.NewRecorder()
	AuthMiddleWare(next, *cnf).ServeHTTP(w, request)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Non expected code: %v", w.Code)
	}
}
