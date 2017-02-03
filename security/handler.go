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
	"fmt"
	"log"
	"net/http"

	"github.com/DamienFontaine/lunarc/web"
	jwt "github.com/dgrijalva/jwt-go"
)

//TokenHandler manage authorizations
func TokenHandler(next http.Handler, cnf web.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				//TODO: On ne passe jamais à l'intérieur
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cnf.Jwt.Key), nil
		})
		if err == nil && token.Valid {
			next.ServeHTTP(w, r)
		} else {
			if r.URL.String() == "/" {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	})
}

//Oauth2 manage authorizations
func Oauth2(next http.Handler, cnf web.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cnf.Jwt.Key), nil
		})
		if err == nil && token.Valid {
			next.ServeHTTP(w, r)
		} else {
			log.Printf("Problem %v", err)
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}
