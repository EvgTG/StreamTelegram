package minidb

import (
	"github.com/recoilme/pudge"
	"github.com/rotisserie/eris"
)

type Pudge struct {
	db *pudge.Db

	videoIDs []string
}

func NewDB() (*Pudge, error) {
	minidb, err := pudge.Open("files/minidb/minidb", &pudge.Config{SyncInterval: 1})
	if err != nil {
		return nil, eris.Wrap(err, "pudge.Open()")
	}

	p := Pudge{
		db: minidb,
	}

	err = p.SetVideoIDs()
	if err != nil {
		return nil, eris.Wrap(err, "p.SetVideoIDs()")
	}

	return &p, nil
}

/*
func (p *Pudge) name() error {
	return err
}
*/
