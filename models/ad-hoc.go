package models

import (
	"database/sql"
	"encoding/json"
	"time"

	uuid "github.com/talkative-ai/go.uuid"
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
	Category    ProjectCategory
	Tags        ProjectTagArray
	ProjectData ProjectItemArray
	TriggerData ProjectTriggerItemArray
}
type ProjectTriggerItem struct {
	ProjectID   uuid.UUID
	ZoneID      uuid.UUID
	TriggerType TriggerType
	RawLBlock
}

type ProjectItemArray []ProjectItem
type ProjectTriggerItemArray []ProjectTriggerItem

func (a *ProjectItemArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}
func (a *ProjectTriggerItemArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}

type ProjectMetadata struct {
	Status      PublishStatus
	Review      *ProjectReviewPublic
	PublishTime time.Time
}
