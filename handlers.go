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
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"

	"github.com/DamienFontaine/lunarc/config"
)

//AuthMiddleWare manage authorizations
func AuthMiddleWare(next http.Handler, cnf config.Server) http.Handler {
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

//Logging logs http requests
func Logging(next http.Handler, log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := StatusResponseWriter{w, 0, 0}
		start := time.Now()
		next.ServeHTTP(&srw, r)
		end := time.Now()
		latency := end.Sub(start)

		log.WithField("client", r.RemoteAddr).WithField("latency", latency).WithField("length", srw.Length()).WithField("code", srw.Status()).Printf("%s %s %s", r.Method, r.URL, r.Proto)
	})
}

//SingleFile returns a handler
func SingleFile(filename string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}
