package models

import (
	"time"

	uuid "github.com/artificial-universe-maker/go.uuid"
	"github.com/go-gorp/gorp"
)

type EventUserActon struct {
	AumModel
	UserID   uuid.UUID
	PubID    uuid.UUID
	RawInput string
}

type EventStateChange struct {
	EventUserActionID string
	StateObject       MutableRuntimeState
	CreatedAt         gorp.NullTime `json:"CreatedAt,omitempty"`
}

func (m *EventStateChange) PreInsert(s gorp.SqlExecutor) error {
	m.CreatedAt.Time = time.Now()
	m.CreatedAt.Valid = true
	return nil
}
