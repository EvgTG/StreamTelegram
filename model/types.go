package model

type Settings struct {
	DBPriority PrioritiesAndVariable
	VideoIDs   []string
}

type PrioritiesAndVariable struct {
	ToID   []int64
	ToIDbl bool
}
