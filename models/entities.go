package models

import (
	"database/sql"
)

type AumModel struct {
	ID      *uint64 `json:"id" db:"id, primarykey, autoincrement"`
	Created *string `json:"created_at,omitempty" db:"-"`
}
type AumProject struct {
	AumModel

	Title     string        `json:"title" db:"title"`
	OwnerID   string        `json:"-" db:"owner_id"`
	StartZone sql.NullInt64 `json:"startZone,omitempty" db:"start_zone_id"` // Expected Zone ID

	Actors []AumActor `json:"actors,omitempty" db:"-"`
	Zones  []AumZone  `json:"locations,omitempty" db:"-"`
	Notes  []AumNote  `json:"notes,omitempty" db:"-"`
}

// AumEntityID
type AEID int

const (
	AEIDActors AEID = iota
)

type AumDialog struct {
	AumModel
	Nodes []AumDialogNode `json:"nodes"`
}

type AumDialogNode struct {
	EntryInput     []AumDialogInput `json:"entry"`
	LogicalSet     RawLBlock        `json:"logical_set"`
	ConnectedNodes []AumDialogNode
}

// Valid dialog types
// Verb statement (Example: “I will <Verb> the <Actor>”)
// Provides a Verb and an Actor
// Verb question (Example: “Did you <Verb> the <Actor>?”)
// Provides a Verb and an Actor
// Possessional question (Example: “Do you have <Actor>?”)
// Provides an Actor
// Greeting (Example: “Hello <Actor>”)
// Provides an Actor
type AumDialogInput string

const (
	AumDialogInputStatementVerb        AumDialogInput = "statement_verb"
	AumDialogInputGreeting             AumDialogInput = "statement_greeting"
	AumDialogInputFarewell             AumDialogInput = "statement_farewell"
	AumDialogInputQuestionVerb         AumDialogInput = "question_verb"
	AumDialogInputQuestionPossessional AumDialogInput = "question_possessional"
	AumDialogInputCustom               AumDialogInput = "custom"
)

type AumActor struct {
	AumModel

	Title   string `json:"title" db:"title"`
	Dialogs []AumDialog
}

type AumZone struct {
	AumModel

	Description      string                `json:"description"`
	Objects          []uint64              `json:"objects,omitempty"`
	Actors           []uint64              `json:"actors,omitempty"`
	LinkedZones      []AumZoneLink         `json:"linkedZones,omitempty"`
	CustomProperties []AumCustomProperties `json:"customProperties,omitempty"`
}

type AumZoneLink struct {
	AumModel

	ZoneFrom uint64
	ZoneTo   uint64
}

type AumTrigger struct {
	AumModel
}

type AumNote struct {
	AumModel
	Text string `json:"text"`
}

type AumCustomProperties map[string]interface{}
