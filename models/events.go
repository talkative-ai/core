package models

import (
	"time"

	"github.com/go-gorp/gorp"
	uuid "github.com/talkative-ai/go.uuid"
)

type EventUserActon struct {
	Model
	UserID   uuid.UUID
	PubID    uuid.UUID
	RawInput string
}

type EventStateChange struct {
	EventUserActionID string
	StateObject       MutableAIRequestState
	CreatedAt         gorp.NullTime `json:"CreatedAt,omitempty"`
}

func (m *EventStateChange) PreInsert(s gorp.SqlExecutor) error {
	m.CreatedAt.Time = time.Now()
	m.CreatedAt.Valid = true
	return nil
}
