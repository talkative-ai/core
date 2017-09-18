package models

import (
	"database/sql"

	"github.com/artificial-universe-maker/go-utilities/common"
)

type ProjectItem struct {
	ProjectID            uint64
	Title                string
	ZoneID               uint64
	ActorID              uint64
	DialogID             uint64
	DialogEntry          common.StringArray
	ParentDialogID       sql.NullInt64
	ChildDialogID        sql.NullInt64
	IsRoot               bool
	LogicalSetAlways     string
	LogicalSetStatements sql.NullString
	LogicalSetID         uint64
	RawLBlock
}
