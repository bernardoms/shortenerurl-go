package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLShortener struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
	OriginalUrl string `json:"originalURL" bson:"originalURL"`
	Alias string `json:"alias" bson:"alias"`
	RedirectCount int `json:"redirectCount" bson:"redirectCount"`
}

type Repository interface {
	Update(urlShortener *URLShortener) (*URLShortener, error)
	Save(urlShortener *URLShortener) (*URLShortener, error)
	FindByAlias(alias string) (*URLShortener, error)
	FindAll() ([] *URLShortener, error)
}