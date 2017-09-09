package models

import (
	"time"

	"github.com/go-gorp/gorp"
)

type EventUserActon struct {
	AumModel
	UserID   uint64
	PubID    uint64
	RawInput string
}

type EventStateChange struct {
	EventUserActionID uint64
	StateObject       MutableRuntimeState
	CreatedAt         gorp.NullTime `json:"CreatedAt,omitempty"`
}

func (m *EventStateChange) PreInsert(s gorp.SqlExecutor) error {
	m.CreatedAt.Time = time.Now()
	m.CreatedAt.Valid = true
	return nil
}
