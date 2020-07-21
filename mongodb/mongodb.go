package mongodb

import (
	"StreamTelegram/go-log"
	"StreamTelegram/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mong struct {
	Mongo   *mongo.Client
	NameCol string
	VIDL    model.VideoIDList
}

func NewDB(nameCol string) *Mong {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Connected to MongoDB!")

	Mong := Mong{client, nameCol, model.VideoIDList{}}

	err = Mong.GetLs(&Mong.VIDL)
	if err != nil {
		log.Fatal(err)
	}

	return &Mong
}

// true  - есть в базе, идем дальше
// false - новый стрим, добавляем в базу
func (m *Mong) Check(id string) (bool, error) {
	for _, a := range m.VIDL.VideoIDs {
		if id == a {
			return true, nil
		}
	}

	if len(m.VIDL.VideoIDs) >= 5 {
		m.VIDL.VideoIDs = append(m.VIDL.VideoIDs[1:5], id)
	} else {
		m.VIDL.VideoIDs = append(m.VIDL.VideoIDs, id)
	}

	err := m.SetLs(&m.VIDL)
	return false, err
}

func (m *Mong) GetLs(ls *model.VideoIDList) error {
	err := m.Mongo.Database(m.NameCol).Collection("VideoIDList").FindOne(context.TODO(), bson.D{{}}).Decode(&ls)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil
		}
	}
	return err
}

func (m *Mong) SetLs(ls *model.VideoIDList) error {
	err := m.Mongo.Database(m.NameCol).Collection("VideoIDList").FindOneAndReplace(context.TODO(), bson.D{{}}, ls).Err()

	if err.Error() == "mongo: no documents in result" {
		_, err = m.Mongo.Database(m.NameCol).Collection("VideoIDList").InsertOne(context.TODO(), ls)
	}

	return err
}
