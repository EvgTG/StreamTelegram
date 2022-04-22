package mongodb

import (
	"StreamTelegram/go-log"
	"StreamTelegram/model"
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mong struct {
	Mongo    *mongo.Client
	NameCol  string
	Settings model.Settings
}

func NewDB(nameCol, mongoUrl string) *Mong {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Fatal(errors.Wrap(err, "mongo.Connect "+mongoUrl))
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "client.Ping"))
	}

	log.Info("Connected to MongoDB!")

	Mong := Mong{client, nameCol, model.Settings{}}

	err = Mong.Mongo.Database(Mong.NameCol).Collection("Settings").FindOne(context.TODO(), bson.D{{}}).Decode(&Mong.Settings)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			log.Fatal(errors.Wrap(err, "Collection(\"Settings\")"))
		}
	}

	return &Mong
}

// true  - есть в базе, идем дальше
// false - новый стрим, добавляем в базу
func (m *Mong) Check(id string) (bool, error) {
	for _, a := range m.Settings.VideoIDs {
		if id == a {
			return true, nil
		}
	}

	if len(m.Settings.VideoIDs) >= 50 {
		m.Settings.VideoIDs = append(m.Settings.VideoIDs[1:50], id)
	} else {
		m.Settings.VideoIDs = append(m.Settings.VideoIDs, id)
	}

	err := m.SetLs(&m.Settings)
	return false, err
}

func (m *Mong) GetLs() model.Settings {
	return m.Settings
}

func (m *Mong) SetLs(ls *model.Settings) error {
	err := m.Mongo.Database(m.NameCol).Collection("Settings").FindOneAndReplace(context.TODO(), bson.D{{}}, ls).Err()

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			_, err = m.Mongo.Database(m.NameCol).Collection("Settings").InsertOne(context.TODO(), ls)
		}
	}

	if err == nil {
		m.Settings = *ls
	}

	return err
}
