package minidb

import (
	"github.com/recoilme/pudge"
)

type Pudge struct {
	db *pudge.Db
}

func NewDB() (*Pudge, error) {
	minidb, err := pudge.Open("files/minidb/minidb", &pudge.Config{SyncInterval: 1})
	if err != nil {
		return nil, err
	}

	p := Pudge{
		db: minidb,
	}

	return &p, nil
}
