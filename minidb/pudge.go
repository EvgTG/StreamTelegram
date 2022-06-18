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

func (p *Pudge) SetCycleDuration(dur int) error {
	err := p.db.Set("dur", dur)
	return err
}

func (p *Pudge) GetCycleDuration() (int, error) {
	dur := 0
	err := p.db.Get("dur", &dur)
	return dur, err
}

func (p *Pudge) SetLocs(locs []string) error {
	err := p.db.Set("locs", locs)
	return err
}

func (p *Pudge) GetLocs() ([]string, error) {
	locs := []string{}
	err := p.db.Get("locs", &locs)
	if err == pudge.ErrKeyNotFound {
		return locs, nil
	}
	return locs, err
}

func (p *Pudge) SetTimeWithCity(bl bool) error {
	err := p.db.Set("timecity", bl)
	return err
}

func (p *Pudge) GetTimeWithCity() (bool, error) {
	bl := true
	err := p.db.Get("timecity", &bl)
	if err == pudge.ErrKeyNotFound {
		return bl, nil
	}
	return bl, err
}

type Channel struct {
	ID          int64
	EndOfStream bool
}

func (p *Pudge) SetNotifyList(list []Channel) error {
	err := p.db.Set("notifylist", list)
	return err
}

func (p *Pudge) GetNotifyList() ([]Channel, error) {
	list := []Channel{}
	err := p.db.Get("notifylist", &list)
	if err == pudge.ErrKeyNotFound {
		return list, nil
	}
	return list, err
}

/*
func (p *Pudge) name() error {
	return err
}
*/
