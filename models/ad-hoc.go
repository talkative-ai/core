package models

import (
	"database/sql"
	"encoding/json"

	uuid "github.com/artificial-universe-maker/go.uuid"
)

type ProjectItemJSONSafe struct {
	ProjectID            string
	Title                string
	ZoneID               string
	ActorID              string
	DialogID             string
	DialogEntry          []string
	ParentDialogID       string
	ChildDialogID        string
	IsRoot               bool
	UnknownHandler       bool
	LogicalSetAlways     string
	LogicalSetStatements string
	RawLBlock
}

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
type VersionedProjectJSONSafe struct {
	ProjectID   uuid.UUID
	Version     int64
	Title       string
	Category    AumProjectCategory
	Tags        AumProjectTagArray
	ProjectData ProjectItemJSONSafeArray
}

type ProjectTriggerItem struct {
	ProjectID   uuid.UUID
	ZoneID      uuid.UUID
	TriggerType AumTriggerType
	RawLBlock
}

type ProjectItemArray []ProjectItem
type ProjectItemJSONSafeArray []ProjectItemJSONSafe

func (a *ProjectItemArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}
