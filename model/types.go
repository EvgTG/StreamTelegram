package model

import "time"

type Settings struct {
	DBPriority PrioritiesAndVariable
	VideoIDs   []string
}

type PrioritiesAndVariable struct {
	ToIDBL      bool
	ToID        []int64
	CycleTimeBL bool
	CycleTime   time.Duration
}
