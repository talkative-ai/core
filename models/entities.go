package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/artificial-universe-maker/core/common"
	uuid "github.com/artificial-universe-maker/go.uuid"
	"github.com/go-gorp/gorp"
)

// AumModel is an embedded struct of common model fields
type AumModel struct {
	ID          uuid.UUID     `db:"ID, primarykey, autoincrement"`
	CreateID    *string       `db:"-" json:",omitempty"`
	CreatedAt   gorp.NullTime `json:"CreatedAt,omitempty"`
	PatchAction *PatchAction  `json:",omitempty" db:"-"`
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
	TeamID      uuid.UUID
	StartZoneID sql.NullString // Expected Zone ID
	IsPrivate   bool

	Actors               []AumActor                `db:"-"`
	Zones                []AumZone                 `db:"-"`
	ZoneActors           []AumZoneActor            `db:"-"`
	PrivateProjectGrants []AumPrivateProjectGrants `db:"-"`
	Notes                []AumNote                 `db:"-"`
}

type AumPublishedProject struct {
	ProjectID uuid.UUID
	TeamID    uuid.UUID
	CreatedAt gorp.NullTime `json:"CreatedAt,omitempty"`
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
		result["StartZoneID"] = p.StartZoneID.String
	}

	return result
}

type AumPrivateProjectGrants struct {
	ID        uuid.UUID `db:"ID, primarykey, autoincrement"`
	ProjectID uuid.UUID
	UserID    uuid.UUID
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

	IsRoot     *bool
	ProjectID  uuid.UUID `json:"-"`
	ActorID    uuid.UUID `json:"-"`
	EntryInput AumDialogInputArray
	RawLBlock
	ChildNodes  *[]*AumDialogNode `db:"-" json:"-"`
	ParentNodes *[]*AumDialogNode `db:"-" json:"-"`
}

func (a *AumDialogNode) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}

// AumDialogInput indicates valid dialog entry types
// As specified in https://aum.ai
// Greeting (Example: “Hello <Actor>”)
// Provides an Actor
type AumDialogInput string
type AumDialogInputArray []AumDialogInput

func (a *AumDialogInputArray) Scan(src interface{}) error {
	arr := common.StringArray{}
	err := arr.Scan(src)
	if err != nil {
		return err
	}
	newA := make(AumDialogInputArray, len(arr.Val))
	for idx, v := range arr.Val {
		newA[idx] = AumDialogInput(v)
	}
	*a = newA
	return nil
}

func (arr *AumDialogInputArray) Value() (driver.Value, error) {
	v := []string{}
	for _, a := range *arr {
		v = append(v, fmt.Sprintf("\"%v\"", string(a)))
	}

	s := strings.Join(v, ",")
	return fmt.Sprintf("{%v}", s), nil
}

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

type PatchAction uint8

const (
	PatchActionCreate PatchAction = iota
	PatchActionDelete
	PatchActionUpdate
)

type UUIDCreateID struct {
	uuid.UUID
	CreateID *string
}

func (u *UUIDCreateID) UnmarshalText(text []byte) error {
	if strings.HasPrefix(string(text), "create") {
		str := string(text)
		u.CreateID = &str
		return nil
	}
	tmp := uuid.UUID{}
	err := tmp.UnmarshalText(text)
	if err != nil {
		return err
	}

	(*u).UUID = tmp
	return nil
}

// AumActor model for the Actor entities
type AumActor struct {
	AumModel

	Title           string
	ProjectID       uuid.UUID           `json:"-"`
	Dialogs         []AumDialogNode     `json:",omitempty" db:"-"`
	DialogRelations []AumDialogRelation `json:",omitempty" db:"-"`
}

type AumDialogRelation struct {
	ParentNodeID UUIDCreateID
	ChildNodeID  UUIDCreateID
	PatchAction  *PatchAction `json:",omitempty" db:"-"`
}

// AumZone model for the Zone entities
type AumZone struct {
	AumModel

	ProjectID   uuid.UUID `json:"-"`
	Title       string
	Description string
	Triggers    map[AumTriggerType]AumTrigger `db:"-"`
}

type AumZoneActor struct {
	ZoneID      UUIDCreateID
	ActorID     UUIDCreateID
	PatchAction *PatchAction `json:",omitempty" db:"-"`
}

type AumTriggerType int

const (
	AumTriggerInitializeZone AumTriggerType = iota
	AumTriggerEnterZone
	AumTriggerExitZone
	AumTriggerVariableUpdate
)

// AumTrigger model for Trigger entities
type AumTrigger struct {
	TriggerType AumTriggerType
	ZoneID      UUIDCreateID
	RawLBlock
	PatchAction *PatchAction `json:",omitempty" db:"-"`
}

// AumNote model for the Note entities
type AumNote struct {
	AumModel
	Text string
}
