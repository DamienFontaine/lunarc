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

package services

import (
	"github.com/DamienFontaine/lunarc/models"
	"github.com/DamienFontaine/lunarc/utils"
	"gopkg.in/mgo.v2/bson"
)

//IUserService interface
type IUserService interface {
	GetByID(id string) models.User
	Get(username string, password string) (models.User, error)
	Add(user models.User) error
	FindAll() []models.User
	Delete(user models.User)
	Update(id string, user models.User) error
}

//UserService works with models.User
type UserService struct {
	MongoService
}

//Get retourne l'utilisateur si celui-ci existe
func (u *UserService) Get(username string, password string) (models.User, error) {
	mongo := u.MongoService.Mongo.Copy()
	defer mongo.Close()

	userCollection := mongo.Database.C("user")
	var user models.User
	userCollection.Find(bson.M{"username": username}).One(&user)

	valid, err := utils.CheckPassword([]byte(password), []byte(user.Salt), []byte(user.Password))
	if err != nil {
		return models.User{}, err
	}
	if valid {
		return user, nil
	}
	return models.User{}, err
}

//GetByID retourne l'utilisateur d'après son ID
func (u *UserService) GetByID(id string) models.User {
	mongo := u.MongoService.Mongo.Copy()
	defer mongo.Close()

	userCollection := mongo.Database.C("user")
	var user models.User
	userCollection.FindId(bson.ObjectIdHex(id)).One(&user)

	return user
}

//FindAll retourne tout les utilisateurs
func (u *UserService) FindAll() []models.User {
	mongo := u.MongoService.Mongo.Copy()
	defer mongo.Close()

	userCollection := mongo.Database.C("user")
	var users []models.User
	userCollection.Find(nil).All(&users)
	return users
}

//Add ajoute un nouvel utilisateur
func (u *UserService) Add(user models.User) error {
	mongo := u.MongoService.Mongo.Copy()
	defer mongo.Close()
	id := bson.NewObjectId()

	salt, err := utils.GenerateSalt()
	if err != nil {
		return err
	}
	user.Salt = string(salt[:32])

	password, err := utils.HashPassword([]byte(user.Password), salt)
	if err != nil {
		return err
	}
	user.Password = string(password[:32])

	userCollection := mongo.Database.C("user")
	userCollection.Insert(&models.User{id, user.Username, user.Firstname, user.Lastname, user.Password, user.Salt, user.Email})

	return nil
}

//Delete supprime un utilisateur
func (u *UserService) Delete(user models.User) {
	mongo := u.MongoService.Mongo.Copy()
	defer mongo.Close()
	userCollection := mongo.Database.C("user")
	userCollection.Remove(bson.M{"_id": user.ID, "username": user.Username})
}

//Update modifie un utilisateur existant
func (u *UserService) Update(id string, user models.User) error {
	mongo := u.MongoService.Mongo.Copy()
	defer mongo.Close()

	salt, err := utils.GenerateSalt()
	if err != nil {
		return err
	}
	user.Salt = string(salt[:32])

	password, err := utils.HashPassword([]byte(user.Password), salt)
	if err != nil {
		return err
	}
	user.Password = string(password[:32])

	userCollection := mongo.Database.C("user")
	userCollection.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"username": user.Username, "lastname": user.Lastname, "firstname": user.Firstname, "password": user.Password, "salt": user.Salt, "email": user.Email}})

	return nil
}
