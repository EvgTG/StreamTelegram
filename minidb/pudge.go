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

func (p *Pudge) SetChannelID(id string) error {
	err := p.db.Set("id", id)
	return err
}

func (p *Pudge) GetChannelID() (string, error) {
	id := ""

	err := p.db.Get("id", &id)
	if err == pudge.ErrKeyNotFound {
		return "", nil
	}

	return id, err
}

/*
func (p *Pudge) name() error {
	return err
}
*/
