package configs

import "os"

type MongoConfig struct {
	MongoURI string
}

func NewMongoConfig()  * MongoConfig{
	return &MongoConfig{
		MongoURI: os.Getenv("mongoURI"),
	}
}