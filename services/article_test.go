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
	"testing"

	"github.com/DamienFontaine/lunarc/datasource"
	"github.com/DamienFontaine/lunarc/models"
	"github.com/DamienFontaine/lunarc/utils"
)

var articleService ArticleService

func BeforeEach() {
	var configUtil = new(utils.ConfigUtil)
	cnf, _ := configUtil.Construct("config.yml", "staging")

	session := datasource.GetSession(cnf.Mongo.Port, cnf.Mongo.Host)
	mongo := datasource.Mongo{Session: session, Database: session.DB(cnf.Mongo.Database)}
	articleService = ArticleService{MongoService: MongoService{Mongo: mongo}}
}

func TestArticleService(t *testing.T) {
	BeforeEach()

	var i interface{} = &articleService
	_, ok := i.(IArticleService)

	if !ok {
		t.Fatalf("ArticleService must implement IArticleService")
	}
}

func TestGetByIdNormal(t *testing.T) {
	BeforeEach()

	article, err := articleService.GetByID("5654921f1d41c84041000001")

	if err != nil {
		t.Fatalf("Mustn't return an error")
	}

	if reflect.DeepEqual(article, models.Article{}) {
		t.Fatalf("Must return an article")
	}
}

func TestGetByIdWithBadIdError(t *testing.T) {
	BeforeEach()

	article, err := articleService.GetByID("5654921f1d41c84041000002")

	if err == nil {
		t.Fatalf("Must return an error")
	}

	if !reflect.DeepEqual(article, models.Article{}) {
		t.Fatalf("Mustn't return an article")
	}
}

func TestGetByPrettyNormal(t *testing.T) {
	BeforeEach()

	article, err := articleService.GetByPretty("first-article")

	if err != nil {
		t.Fatalf("Mustn't return an error")
	}

	if reflect.DeepEqual(article, models.Article{}) {
		t.Fatalf("Mustn't return an article")
	}
}

func TestGetByPrettyWithBadPrettyError(t *testing.T) {
	BeforeEach()

	article, err := articleService.GetByPretty("first-articl")

	if err == nil {
		t.Fatalf("Must return an error")
	}

	if !reflect.DeepEqual(article, models.Article{}) {
		t.Fatalf("Must return an article")
	}
}

func TestFindByStatusNormal(t *testing.T) {
	BeforeEach()

	articles, err := articleService.FindByStatus("PUBLISH")

	if err != nil {
		t.Fatalf("Mustn't return an error")
	}

	if len(articles) != 2 {
		t.Fatalf("Must return 2 articles but %d returned", len(articles))
	}
}

func TestFindByStatusWithBadStatusError(t *testing.T) {
	BeforeEach()

	articles, err := articleService.FindByStatus("BAD")

	if err != nil {
		t.Fatalf("Must return an error")
	}

	if len(articles) != 0 {
		t.Fatalf("Must return 0 articles but %d returned", len(articles))
	}
}

func TestFindAllNormal(t *testing.T) {
	BeforeEach()

	articles, err := articleService.FindAll()

	if err != nil {
		t.Fatalf("Mustn't return an error")
	}

	if len(articles) != 4 {
		t.Fatalf("Must return 4 articles but %d returned", len(articles))
	}
}

func TestAddNormal(t *testing.T) {
	BeforeEach()

	article := models.Article{Titre: "New Article"}

	article, err := articleService.Add(article)

	if err != nil {
		t.Fatalf("Mustn't return an error")
	}

	if reflect.DeepEqual(article, models.Article{}) {
		t.Fatalf("Must return an article")
	}
	articleService.Delete(article)
}

func TestUpdateNormal(t *testing.T) {
	BeforeEach()

	oldTitle := "New Article 2"
	newTitle := "New Article 2.1"
	pretty := utils.SanitizeTitle(newTitle)
	article := models.Article{Titre: oldTitle}
	article, _ = articleService.Add(article)

	article.Titre = newTitle
	err := articleService.Update(string(article.ID.Hex()), article)

	if err != nil {
		t.Fatalf("Mustn't return error: %s", err)
	}

	article, err = articleService.GetByPretty(pretty)

	if err != nil {
		t.Fatalf("Mustn't return error: %s", err)
	}

	if article.Titre != newTitle {
		t.Fatalf("Must update title %s to %s but %s return", oldTitle, newTitle, article.Titre)
	}

	articleService.Delete(article)
}

func TestUpdateBadIdError(t *testing.T) {
	BeforeEach()

	article := models.Article{Titre: "Le titre"}

	err := articleService.Update("5654921f1d41c84041000002", article)

	if err == nil {
		t.Fatalf("Must return error")
	}
}

func TestDeleteNormal(t *testing.T) {
	BeforeEach()

	article := models.Article{Titre: "New Article Delete"}
	article, _ = articleService.Add(article)

	err := articleService.Delete(article)

	if err != nil {
		t.Fatalf("Mustn't return error: %s", err)
	}

	article, err = articleService.GetByPretty(article.Pretty)

	if err == nil {
		t.Fatalf("Must return an error")
	}
}

func TestDeleteBadIdError(t *testing.T) {
	BeforeEach()

	article := models.Article{Titre: "Le titre"}

	err := articleService.Delete(article)

	if err == nil {
		t.Fatalf("Must return error")
	}
}
