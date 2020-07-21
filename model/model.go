package model

type db interface {
	GetLs(ls *VideoIDList) error
	SetLs(ls *VideoIDList) error
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

func (m *Model) GetLs(ls *VideoIDList) error {
	return m.db.GetLs(ls)
}

func (m *Model) SetLs(ls *VideoIDList) error {
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
