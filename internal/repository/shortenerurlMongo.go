package repository

import (
	"context"
	"github.com/bernardoms/shortenerurl-go/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)
type SessionCreator struct {
	Client *mongo.Client
}

type Mongo struct {
	Collection *mongo.Collection
}

var session *mongo.Client

func New(config *configs.MongoConfig) {

	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))

	m := SessionCreator	{
		Client: client,
	}

	ctx := context.TODO()

	if err != nil {
		log.Print("Error on connecting to database ", err)
	}

	err = m.Client.Connect(ctx)

	session = client

	if err != nil {
		log.Print("Error on connecting to database ", err)
	}
}

func GetShortenerCollection () *mongo.Collection {
	c := session.Database("test").Collection("shorteners")
	return c
}

func (m Mongo) FindAll() ([] *URLShortener, error) {
	var results []*URLShortener

	cur ,err := m.Collection.Find(context.TODO(), bson.D{},
	options.Find().SetSort(bson.D{{"redirectCount", -1}}))

	if err == nil {

		for cur.Next(context.TODO()) {
			var elem URLShortener
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			results = append(results, &elem)
		}
	}
	return results, err
}

func (m Mongo) Save(urlShortener *URLShortener) (*URLShortener, error){

	_ , err := m.Collection.InsertOne(context.TODO(), &urlShortener)

	if err != nil{
		log.Printf("Error inserting %o with error %s", &urlShortener, err)
	}

	return urlShortener, err
}

func (m Mongo) 	FindByAlias(alias string) (*URLShortener, error) {
	var result *URLShortener

	s := m.Collection.FindOne(context.TODO(), bson.M{"alias": alias})

	err := s.Decode(&result)

	if err != nil {
		log.Print(err)
	}

	return result, err
}

func (m Mongo) Update(urlShortener *URLShortener) (*URLShortener, error) {

	_, err := m.Collection.UpdateOne(context.TODO(), bson.M{"_id" : urlShortener.Id},
		bson.M{"$set": bson.M{"redirectCount": urlShortener.RedirectCount,
			"originalURL": urlShortener.OriginalUrl,
		"alias" : urlShortener.Alias}})

	if err != nil {
		log.Printf("error updating %s", err)
	}

	return urlShortener, err
}