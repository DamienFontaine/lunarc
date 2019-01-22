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

package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"
)

//Mongo is a datasource.
type Mongo struct {
	Client   *mongo.Client
	Database *mongo.Database
	context  context.Context
}

//NewMongo creates a newinstance of Mongo
func NewMongo(filename string, environment string) (*Mongo, error) {
	ctx := context.Background()
	cnf, err := GetMongo(filename, environment)
	if err != nil {
		return nil, err
	}
	var uri string
	if len(cnf.Username) > 0 && len(cnf.Password) > 0 {
		uri = fmt.Sprintf(`mongodb://%s:%s@%s:%d/%s`,
			cnf.Username,
			cnf.Password,
			cnf.Host,
			cnf.Port,
			cnf.Database,
		)
	} else {
		uri = fmt.Sprintf(`mongodb://%s:%d/%s`,
			cnf.Host,
			cnf.Port,
			cnf.Database,
		)
	}
	client, err := mongo.NewClient(uri)
	if err != nil {
		log.Printf("L'URI du serveur MongoDB est incorrect: %s", uri)
		return nil, err
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Print("Impossible d'utiliser ce context")
		return nil, err
	}

	db := client.Database(cnf.Database)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Impossible de contacter %v sur le port %d", cnf.Host, cnf.Port)
		return nil, err
	}
	return &Mongo{Client: client, Database: db, context: ctx}, nil
}

//Disconnect a Mongo client
func (m *Mongo) Disconnect() error {
	err := m.Client.Disconnect(m.context)
	if err != nil {
		log.Printf("Impossible de fermer la connexion")
		return err
	}
	return nil
}
