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
	"errors"

	"github.com/DamienFontaine/lunarc/models"
	"github.com/DamienFontaine/lunarc/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//IArticleService interface
type IArticleService interface {
	GetByID(id string) (models.Article, error)
	GetByPretty(pretty string) (models.Article, error)
	FindByStatus(status string) ([]models.Article, error)
	Add(article models.Article) (models.Article, error)
	FindAll() ([]models.Article, error)
	Delete(article models.Article) error
	Update(id string, article models.Article) error
}

//ArticleService works with models.Article
type ArticleService struct {
	MongoService
}

//GetByID retourne l'article d'après son ID
func (a *ArticleService) GetByID(id string) (article models.Article, err error) {
	mongo := a.MongoService.Mongo.Copy()
	defer mongo.Close()

	articleCollection := mongo.Database.C("article")
	err = articleCollection.FindId(bson.ObjectIdHex(id)).One(&article)

	if err != nil {
		return article, errors.New("No article")
	}

	return article, nil
}

//GetByPretty retourne l'article d'après son Pretty
func (a *ArticleService) GetByPretty(pretty string) (article models.Article, err error) {
	mongo := a.MongoService.Mongo.Copy()
	defer mongo.Close()

	articleCollection := mongo.Database.C("article")
	err = articleCollection.Find(bson.M{"pretty": pretty}).One(&article)

	if err != nil {
		return article, errors.New("No article")
	}

	return article, nil
}

//Add ajoute un nouvel article
func (a *ArticleService) Add(article models.Article) (models.Article, error) {
	mongo := a.MongoService.Mongo.Copy()
	defer mongo.Close()
	id := bson.NewObjectId()
	pretty := utils.SanitizeTitle(article.Titre)
	articleCollection := mongo.Database.C("article")
	articleCollection.Insert(&models.Article{id, article.Titre, pretty, article.Texte, article.Tags, article.Image, article.Vignette, article.Status, article.Create, article.Create, mgo.DBRef{Collection: "user", Id: article.UserRef.Id}})

	err := articleCollection.FindId(id).One(&article)

	if err != nil {
		return models.Article{}, err
	}

	return article, nil
}

//FindByStatus retourne les articles d'après leur status
func (a *ArticleService) FindByStatus(status string) (articles []models.Article, err error) {
	mongo := a.MongoService.Mongo.Copy()
	defer mongo.Close()

	articleCollection := mongo.Database.C("article")
	err = articleCollection.Find(bson.M{"status": status}).All(&articles)

	if err != nil {
		return articles, errors.New("Error in FindByStatus")
	}

	return articles, nil
}

//FindAll retourne tout les articles
func (a *ArticleService) FindAll() (articles []models.Article, err error) {
	mongo := a.MongoService.Mongo.Copy()
	defer mongo.Close()

	articleCollection := mongo.Database.C("article")
	err = articleCollection.Find(nil).All(&articles)

	if err != nil {
		return articles, errors.New("Error in FindAll")
	}

	return articles, nil
}

//Delete supprime un article
func (a *ArticleService) Delete(article models.Article) (err error) {
	mongo := a.MongoService.Mongo.Copy()
	defer mongo.Close()
	articleCollection := mongo.Database.C("article")
	err = articleCollection.Remove(bson.M{"_id": article.ID, "titre": article.Titre})
	return
}

//Update modifie un article existant
func (a *ArticleService) Update(id string, article models.Article) (err error) {
	mongo := a.MongoService.Mongo.Copy()
	defer mongo.Close()
	pretty := utils.SanitizeTitle(article.Titre)
	articleCollection := mongo.Database.C("article")
	err = articleCollection.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"titre": article.Titre, "pretty": pretty, "image": article.Image, "vignette": article.Vignette, "texte": article.Texte, "status": article.Status, "modified": article.Modified, "tags": article.Tags, "userref": bson.M{"$ref": article.UserRef.Collection, "$id": article.UserRef.Id}}})
	return
}
