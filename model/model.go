package model

type db interface {
	GetLs() Settings
	SetLs(ls *Settings) error
	Check(id string) (bool, error)
}

type Model struct {
	db
}

func New(db db) *Model {
	return &Model{
		db: db,
	}
}

func (m *Model) GetLs() Settings {
	return m.db.GetLs()
}

func (m *Model) SetLs(ls *Settings) error {
	return m.db.SetLs(ls)
}

func (m *Model) Check(id string) (bool, error) {
	return m.db.Check(id)
}

/*
func (m *Model) {
	return m.db.()
}
*/
