package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/artificial-universe-maker/go-utilities/common"
	"github.com/go-gorp/gorp"
)

// AumModel is an embedded struct of common model fields
type AumModel struct {
	ID        uint64        `db:"ID, primarykey, autoincrement"`
	CreateID  *int          `db:"-" json:",omitempty"`
	CreatedAt gorp.NullTime `json:"CreatedAt,omitempty"`
}

func (m *AumModel) PreInsert(s gorp.SqlExecutor) error {
	m.CreatedAt.Time = time.Now()
	m.CreatedAt.Valid = true
	return nil
}

// AumProject is the model for a Workbench project
type AumProject struct {
	AumModel

	Title       string
	TeamID      uint64
	StartZoneID sql.NullInt64 // Expected Zone ID

	Actors     []AumActor     `db:"-"`
	Zones      []AumZone      `db:"-"`
	ZoneActors []AumZoneActor `db:"-"`
	Notes      []AumNote      `db:"-"`
}

func (p AumProject) PrepareMarshal() map[string]interface{} {
	result := map[string]interface{}{
		"ID":         p.ID,
		"Title":      p.Title,
		"CreatedAt":  p.CreatedAt.Time,
		"Zones":      p.Zones,
		"Actors":     p.Actors,
		"ZoneActors": p.ZoneActors,
	}

	if p.StartZoneID.Valid {
		result["StartZoneID"] = p.StartZoneID.Int64
	}

	return result
}

// AEID is an AumEntityID
// Useful for Redis key mapping
type AEID int

const (
	// AEIDActor AumEntityID for Actor
	AEIDActor AEID = iota

	// AEIDZone AumEntityID for Zone
	AEIDZone

	// AEIDTrigger AumEntityID for Trigger
	AEIDTrigger
	// AEIDDialogNode AumEntityID for DialogNode
	AEIDDialogNode
	// AEIDActionBundle AumEntityID for ActionBundle
	AEIDActionBundle
)

// AumDialogNode is a single instance of a Dialog
// If ParentNodes == nil, then it's the entry of a dialog
// If ChildNodes == nil, then it's the end of the dialog
type AumDialogNode struct {
	AumModel

	ProjectID   uint64
	ActorID     uint64
	EntryInput  []AumDialogInput
	LogicalSet  RawLBlock
	ChildNodes  *[]*AumDialogNode `db:"-"`
	ParentNodes *[]*AumDialogNode `db:"-"`
}

// AumDialogInput indicates valid dialog entry types
// As specified in https://aum.ai
// Greeting (Example: “Hello <Actor>”)
// Provides an Actor
type AumDialogInput string

const (
	// AumDialogInputStatementVerb Verb statement
	// (Example: “I will <Verb> the <Actor>”)
	// Provides a Verb and an Actor
	AumDialogInputStatementVerb AumDialogInput = "statement_verb"
	// AumDialogInputGreeting Generic greeting
	// (Example: "Hello <Actor>")
	// Provides an optional Actor
	AumDialogInputGreeting AumDialogInput = "statement_greeting"
	// AumDialogInputFarewell Generic farewell
	// (Example: "Goodbye <Actor>")
	// Provides an optional Actor
	AumDialogInputFarewell AumDialogInput = "statement_farewell"
	// AumDialogInputQuestionVerb Verb question
	// (Example: “Did you <Verb> the <Actor>?”)
	// Provides a Verb and an Actor
	AumDialogInputQuestionVerb AumDialogInput = "question_verb"
	// AumDialogInputQuestionPossessional Possessional question
	// (Example: “Do you have <Actor>?”)
	// Provides an Actor
	AumDialogInputQuestionPossessional AumDialogInput = "question_possessional"
)

// AumActor model for the Actor entities
type AumActor struct {
	AumModel

	Title     string
	ProjectID uint64                 `json:"-"`
	ZoneID    *uint64                `json:",omitempty" db:"-"`
	Dialogs   []AumMinimalDialogNode `json:",omitempty" db:"-"`
}

type AumMinimalDialogNode struct {
	ID           string
	LogicalSetID uint64
	Entry        common.StringArray
	Always       AlwaysMap
}

type AlwaysMap map[string][]map[string]interface{}

func (arr *AlwaysMap) Value() (driver.Value, error) {
	return *arr, nil
}

func (arr *AlwaysMap) Scan(src interface{}) error {
	json.Unmarshal(src.([]byte), &arr)
	return nil
}

// AumZone model for the Zone entities
type AumZone struct {
	AumModel

	ProjectID   uint64 `json:"-"`
	Title       string
	Description string
	LinkedZones []AumZoneLink `json:"LinkedZones,omitempty"`
}

type AumZoneActor struct {
	ZoneID  uint64
	ActorID uint64
}

// AumZoneLink explicitly linked zones
type AumZoneLink struct {
	AumModel

	ZoneFrom uint64
	ZoneTo   uint64
}

// AumTrigger model for Trigger entities
type AumTrigger struct {
	AumModel
}

// AumNote model for the Note entities
type AumNote struct {
	AumModel
	Text string
}
