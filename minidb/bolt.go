package minidb

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/rotisserie/eris"
	bolt "go.etcd.io/bbolt"
)

type MiniDB struct {
	bolt *bolt.DB

	videoIDs []string
}

func NewDB() (*MiniDB, error) {
	db, err := bolt.Open("files/my.db", 0666, nil)
	if err != nil {
		return nil, eris.Wrap(err, "bolt.Open()")
	}

	mini := MiniDB{bolt: db}
	mini.bolt.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("my"))
		return err
	})

	err = mini.GetVideoIDs()
	if err != nil && err != ErrKeyNotFound {
		return nil, eris.Wrap(err, "mini.GetVideoIDs()")
	}

	return &mini, nil
}

func (mini *MiniDB) write(key string, obj any) error {
	return mini.bolt.Update(func(tx *bolt.Tx) error {
		bts, err := encode(obj)
		if err != nil {
			return err
		}

		return tx.Bucket([]byte("my")).Put([]byte(key), bts)
	})
}

func (mini *MiniDB) read(key string, obj any) error {
	return mini.bolt.View(func(tx *bolt.Tx) error {
		return decode(tx.Bucket([]byte("my")).Get([]byte(key)), obj)
	})
}

func encode(obj any) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var ErrKeyNotFound = errors.New("Error: key not found")

func decode(s []byte, obj any) error {
	if s == nil {
		return nil
		//return ErrKeyNotFound
	}

	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}

/*
func (mini *MiniDB) name() error {
	return nil
}
*/
