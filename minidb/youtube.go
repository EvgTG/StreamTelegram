package minidb

import (
	"github.com/recoilme/pudge"
	"github.com/rotisserie/eris"
)

// true  - есть в базе, идем дальше
// false - новый стрим, добавляем в базу
func (p *Pudge) CheckVideo(id string) (bool, error) {
	for _, a := range p.videoIDs {
		if id == a {
			return true, nil
		}
	}

	if len(p.videoIDs) >= 50 {
		p.videoIDs = append(p.videoIDs[1:50], id)
	} else {
		p.videoIDs = append(p.videoIDs, id)
	}

	return false, p.SetVideoIDs()
}

func (p *Pudge) SetVideoIDs() error {
	err := p.db.Set("vids", &p.videoIDs)
	if err == pudge.ErrKeyNotFound {
		return eris.Wrap(err, "db.Set(vids)")
	}
	return nil
}

func (p *Pudge) GetVideoIDs() error {
	err := p.db.Get("vids", &p.videoIDs)
	if err == pudge.ErrKeyNotFound {
		return eris.Wrap(err, "db.Get(vids)")
	}
	return nil
}
