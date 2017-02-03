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
package security_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.lineolia.net/meda/leode-server/oauth2"

	"github.com/DamienFontaine/lunarc/mock"
	"github.com/DamienFontaine/lunarc/security"
	"github.com/DamienFontaine/lunarc/web"
	"github.com/golang/mock/gomock"
)

func TestAuthenticateNormal(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	server, _ := web.NewWebServer("config.yml", "test")
	username := "admin"
	password := "admin"
	user := security.User{Username: "admin", Password: "admin", Salt: "salt", Email: "admin@lineolia.net"}
	mockUserManager := mock.NewMockUserManager(mockCtrl)
	authController := security.NewAuthController(mockUserManager, server.GetConfig())
	var jsonStr = []byte(`{"username": "admin", "password": "admin"}`)
	r, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	mockUserManager.EXPECT().Get(username, password).Return(user, nil)
	authController.Authenticate(w, r)
	if !strings.Contains(w.Body.String(), "id_token") {
		t.Fatalf("Non expected Body")
	}
}

func TestAuthenticateBadPassword(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	username := "admin"
	password := "admin"
	user := security.User{}
	server, _ := web.NewWebServer("config.yml", "test")
	mockUserManager := mock.NewMockUserManager(mockCtrl)
	authController := security.NewAuthController(mockUserManager, server.GetConfig())
	var jsonStr = []byte(`{"username": "admin", "password": "admin"}`)
	r, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	mockUserManager.EXPECT().Get(username, password).Return(user, nil)
	authController.Authenticate(w, r)
	if strings.Contains(w.Body.String(), "id_token") {
		t.Fatalf("Non expected Body")
	}
}

func TestAuthenticateBadSignedString(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	server, _ := web.NewWebServer("config.yml", "test")
	mockUserManager := mock.NewMockUserManager(mockCtrl)
	authController := security.NewAuthController(mockUserManager, server.GetConfig())
	var jsonStr = []byte(`{"username": "admin", "password: "admin"}`)
	r, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	authController.Authenticate(w, r)
	if w.Code != 400 {
		t.Fatalf("Non expected return code %v != 400", w.Code)
	}
}

//TestAuthorizeNormal
func TestAuthorizeNormal(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	server, _ := web.NewWebServer("config.yml", "test")

	mockApplicationManager := mock.NewMockApplicationManager(mockCtrl)
	oAuth2Controller := security.NewOAuth2Controller(mockApplicationManager, server.GetConfig())

	r, _ := http.NewRequest("GET", "/oauth2/authorize?client_id=1&response_type=code&redirect_uri=http://redirect", nil)
	w := httptest.NewRecorder()
	oAuth2Controller.Authorize(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Code must be 200 but get %v", w.Code)
	}
}

// TestTokenNormal
func TestTokenNormal(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	server, _ := web.NewWebServer("config.yml", "test")
	mockApplicationManager := mock.NewMockApplicationManager(mockCtrl)
	oAuth2Controller := security.NewOAuth2Controller(mockApplicationManager, server.GetConfig())
	clientID := "1"
	redirectURI := "http://redirect"
	userID := "1"
	sharedKey := "LunarcSecretKey"
	code, _ := oauth2.EncodeOAuth2Code(clientID, redirectURI, userID, sharedKey)
	r, _ := http.NewRequest("POST", fmt.Sprintf("/oauth2/token?grant_type=authorization_code&code=%v", code), nil)
	w := httptest.NewRecorder()
	oAuth2Controller.Token(w, r)
	if strings.Compare(w.Body.String(), "") == 0 {
		t.Fatal("Must return a token")
	}
}
