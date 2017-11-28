package models

import (
	"database/sql"

	"github.com/artificial-universe-maker/core/common"
	uuid "github.com/artificial-universe-maker/go.uuid"
)

type ProjectItem struct {
	ProjectID            uuid.UUID
	Title                string
	ZoneID               uuid.UUID
	ActorID              uuid.UUID
	DialogID             uuid.UUID
	DialogEntry          common.StringArray
	ParentDialogID       uuid.NullUUID
	ChildDialogID        uuid.NullUUID
	IsRoot               bool
	LogicalSetAlways     string
	LogicalSetStatements sql.NullString
	RawLBlock
}

type ProjectTriggerItem struct {
	ProjectID   uuid.UUID
	ZoneID      uuid.UUID
	TriggerType AumTriggerType
	RawLBlock
}
