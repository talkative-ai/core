package models

import (
	"database/sql"
	"encoding/json"

	uuid "github.com/artificial-universe-maker/go.uuid"
)

type ProjectItem struct {
	ProjectID            uuid.UUID
	Title                string
	ZoneID               uuid.UUID
	ActorID              uuid.UUID
	DialogID             uuid.UUID
	DialogEntry          []string
	ParentDialogID       uuid.NullUUID
	ChildDialogID        uuid.NullUUID
	IsRoot               bool
	UnknownHandler       bool
	LogicalSetAlways     string
	LogicalSetStatements sql.NullString
	RawLBlock
}

type VersionedProject struct {
	ProjectID   uuid.UUID
	Version     int64
	Title       string
	Category    AumProjectCategory
	Tags        AumProjectTagArray
	ProjectData ProjectItemArray
}

type ProjectTriggerItem struct {
	ProjectID   uuid.UUID
	ZoneID      uuid.UUID
	TriggerType AumTriggerType
	RawLBlock
}

type ProjectItemArray []ProjectItem

func (a *ProjectItemArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}
