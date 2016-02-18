// +build integration

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
	"reflect"
	"strings"
	"testing"

	"github.com/DamienFontaine/lunarc/datasource"
	"github.com/DamienFontaine/lunarc/models"
)

var userService UserService

func UserBeforeEach() {
	mongo, _ := datasource.NewMongo("config.yml", "staging")
	userService = UserService{MongoService: MongoService{Mongo: *mongo}}
}

func TestUserService(t *testing.T) {
	BeforeEach()

	var i interface{} = &userService
	_, ok := i.(IUserService)

	if !ok {
		t.Fatalf("UserService must implement IUserService")
	}
}

func TestGetNormal(t *testing.T) {
	UserBeforeEach()

	user, err := userService.Get("administrator", "administrator")

	if err != nil {
		t.Fatalf("Mustn't return an error")
	}

	if reflect.DeepEqual(user, models.User{}) {
		t.Fatalf("Must return an article")
	}
}

func TestGetInvalidPasswordError(t *testing.T) {
	UserBeforeEach()

	user, err := userService.Get("administrator", "admin")

	if err == nil {
		t.Fatalf("Must return an error")
	}

	if !reflect.DeepEqual(user, models.User{}) {
		t.Fatalf("Mustn't return an article")
	}
}

func TestGetInvalidUserError(t *testing.T) {
	UserBeforeEach()

	user, err := userService.Get("admin", "administrator")

	if err == nil {
		t.Fatalf("Must return an error")
	}

	if !reflect.DeepEqual(user, models.User{}) {
		t.Fatalf("Mustn't return an article")
	}
}

func TestGetByIDNormal(t *testing.T) {
	UserBeforeEach()

	user, err := userService.GetByID("56781c0e1d41c8e862787d1c")

	if err != nil {
		t.Fatalf("Mustn't return an error")
	}

	if reflect.DeepEqual(user, models.User{}) {
		t.Fatalf("Must return an article")
	}
}

func TestGetByIDBadIDError(t *testing.T) {
	UserBeforeEach()

	user, err := userService.GetByID("56781c0e1d41c8e862787d1")

	if err == nil {
		t.Fatalf("Must return an error")
	}

	if !reflect.DeepEqual(user, models.User{}) {
		t.Fatalf("Mustn't return an article")
	}
}

func TestGetByIDNonExistantIDError(t *testing.T) {
	UserBeforeEach()

	user, err := userService.GetByID("56781c0e1d41c8e862787d12")

	if err == nil {
		t.Fatalf("Must return an error")
	}

	if !reflect.DeepEqual(user, models.User{}) {
		t.Fatalf("Mustn't return an article")
	}
}

func TestUserFindAllNormal(t *testing.T) {
	UserBeforeEach()

	users, err := userService.FindAll()

	if err != nil {
		t.Fatalf("Mustn't return an error")
	}

	if len(users) != 1 {
		t.Fatalf("Expected 1 but %d", len(users))
	}
}

func TestUserAddNormal(t *testing.T) {
	BeforeEach()

	user := models.User{Username: "Test"}

	user, err := userService.Add(user)

	if err != nil {
		t.Fatalf("Mustn't return an error")
	}

	if reflect.DeepEqual(user, models.User{}) {
		t.Fatalf("Must return an user")
	}
	userService.Delete(user)
}

func TestUserUpdateNormal(t *testing.T) {
	BeforeEach()

	oldUsername := "NewUserTest"
	newUsername := "NewUserTest2"
	user := models.User{Username: oldUsername}
	user, err := userService.Add(user)

	user.Username = newUsername
	err = userService.Update(string(user.ID.Hex()), user)

	if err != nil {
		t.Fatalf("Mustn't return error: %s", err)
	}

	user, err = userService.GetByID(user.ID.Hex())

	if strings.Compare(user.Username, newUsername) != 0 {
		t.Fatalf("Expected %s to be %s", user.Username, newUsername)
	}

	if err != nil {
		t.Fatalf("Mustn't return an error")
	}

	if reflect.DeepEqual(user, models.User{}) {
		t.Fatalf("Must return an user")
	}
	userService.Delete(user)
}

func TestUserUpdateBadIdError(t *testing.T) {
	BeforeEach()

	user := models.User{Username: "NewUserTest3"}

	err := userService.Update("5654921f1d41c84041000002", user)

	if err == nil {
		t.Fatalf("Must return an error")
	}
}

func TestUserUpdateIncorrectIdError(t *testing.T) {
	BeforeEach()

	user := models.User{Username: "NewUserTest3"}

	err := userService.Update("5654921f1d41c8404100000", user)

	if err == nil {
		t.Fatalf("Must return an error")
	}
}

func TestUserDeleteNormal(t *testing.T) {
	BeforeEach()

	user := models.User{Username: "NewUserTest4"}
	user, err := userService.Add(user)

	err = userService.Delete(user)

	if err != nil {
		t.Fatalf("Mustn't return error: %s", err)
	}

	user, err = userService.GetByID(string(user.ID.Hex()))

	if err == nil {
		t.Fatalf("Must return an error")
	}
}

func TestUserDeleteBadIdError(t *testing.T) {
	BeforeEach()

	user := models.User{Username: "NewUserTest5"}

	err := userService.Delete(user)

	if err == nil {
		t.Fatalf("Must return error")
	}
}
