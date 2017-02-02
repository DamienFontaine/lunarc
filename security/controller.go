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
	"fmt"
	"net/http"
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
