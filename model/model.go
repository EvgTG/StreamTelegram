package model

type db interface {
}

type Model struct {
	db
}

func New(db db) *Model {
	return &Model{
		db: db,
	}
}

/*
func (m *Model) {
	return m.db.()
}
*/
