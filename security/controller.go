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
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DamienFontaine/lunarc/web"
	jwt "github.com/dgrijalva/jwt-go"
)

//AuthController manages Authentication in Server.
type AuthController struct {
	cnf         web.Config
	UserManager UserManager
}

//NewAuthController constructs new AuthController
func NewAuthController(um UserManager, cnf web.Config) *AuthController {
	authController := AuthController{UserManager: um, cnf: cnf}
	return &authController
}

//Authenticate controls authorizations
func (c *AuthController) Authenticate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user User
	var data map[string]string
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 400)
		return
	}
	user, _ = c.UserManager.Get(user.Username, user.Password)
	if user.Username != "" {
		token := jwt.New(jwt.GetSigningMethod("HS256"))
		token.Claims["username"] = user.Username
		token.Claims["email"] = user.Email
		token.Claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
		tokenString, _ := token.SignedString([]byte(c.cnf.Jwt.Key))
		data = map[string]string{
			"id_token": tokenString,
		}
	}
	js, _ := json.Marshal(data)
	w.Write(js)
}

//OAuth2Controller manages Authentication in Server.
type OAuth2Controller struct {
	cnf                web.Config
	ApplicationManager ApplicationManager
}

//NewOAuth2Controller constructs new AuthController
func NewOAuth2Controller(am ApplicationManager, cnf web.Config) *OAuth2Controller {
	oAuth2Controller := OAuth2Controller{cnf: cnf, ApplicationManager: am}
	return &oAuth2Controller
}

//Token returns a token
func (c *OAuth2Controller) Token(w http.ResponseWriter, r *http.Request) {
	grantType := r.URL.Query().Get("grant_type")
	code := r.URL.Query().Get("code")
	if strings.Compare(grantType, "authorization_code") != 0 {
		http.Error(w, errors.New("Parameter grant_type is required").Error(), http.StatusBadRequest)
		return
	}
	if strings.Compare(code, "") == 0 {
		http.Error(w, errors.New("Parameter code is required").Error(), http.StatusBadRequest)
		return
	}
	response, err := DecodeOAuth2Code(code, c.cnf.Jwt.Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	i, err := strconv.ParseInt(response.Exp, 10, 64)
	exp := time.Unix(i, 0)
	if exp.After(time.Now()) {
		log.Printf("Code is expired")
	} else {
		token := jwt.New(jwt.GetSigningMethod("HS256"))
		token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
		tokenString, _ := token.SignedString([]byte(c.cnf.Jwt.Key))
		data := map[string]string{
			"access_token": tokenString,
			"token_type":   "bearer",
		}
		js, _ := json.Marshal(data)
		w.Write(js)
	}
}

//Authorize user
func (c *OAuth2Controller) Authorize(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	responsetype := r.URL.Query().Get("response_type")
	redirectURI := r.URL.Query().Get("redirect_uri")
	userID := r.URL.Query().Get("user_id")
	if len(redirectURI) == 0 {
		app, _ := c.ApplicationManager.GetByClientID(clientID)
		redirectURI = app.Callback
	}
	if len(redirectURI) == 0 {
		log.Print("Pas de paramètre redirect_uri")
		// return an error code ?
	} else {
		if strings.Compare(responsetype, "code") == 0 {
			code, err := EncodeOAuth2Code(clientID, redirectURI, userID, c.cnf.Jwt.Key)
			if err != nil {
				log.Printf("Error: %v", err)
			}
			data := map[string]string{
				"redirectURI": redirectURI,
				"code":        code,
			}
			js, _ := json.Marshal(data)
			w.Write(js)
		} else {
			log.Print("Pas de paramètre code")
			//return a Token
		}
	}
}
