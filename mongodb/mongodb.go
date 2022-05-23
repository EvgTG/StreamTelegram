package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"streamtg/go-log"
	"sync"
)

type Mong struct {
	Mongo *mongo.Client
	DB    *mongo.Database
	mutex sync.Mutex
}

func NewDB(nameDB, mongoUrl string) *Mong {
	if mongoUrl == "" {
		mongoUrl = "mongodb://127.0.0.1:27017"
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Connected to MongoDB!")

	Mong := Mong{client, client.Database(nameDB), sync.Mutex{}}

	return &Mong
}
