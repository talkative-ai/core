package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/talkative-ai/core/common"
	uuid "github.com/talkative-ai/go.uuid"
)

// Model is an embedded struct of common model fields
type Model struct {
	ID          uuid.UUID     `db:"ID, primarykey, autoincrement"`
	CreateID    *string       `db:"-" json:",omitempty"`
	CreatedAt   gorp.NullTime `json:"CreatedAt,omitempty"`
	PatchAction *PatchAction  `json:",omitempty" db:"-"`
}

func (m *Model) PreInsert(s gorp.SqlExecutor) error {
	m.CreatedAt.Time = time.Now()
	m.CreatedAt.Valid = true
	return nil
}

type ProjectCategory string

const (
	ProjectCategoryEntertainment ProjectCategory = "Entertainment"
	ProjectCategoryMiscellanious                 = "Miscellaneous"
	ProjectCategoryBusiness                      = "Business"
	ProjectCategoryEducation                     = "Education"
)

type ProjectTag string
type ProjectTagArray []ProjectTag

const (
	ProjectTagInteractiveStory ProjectTag = "Interactive Story"
	ProjectTagHumor                       = "Humor"
	ProjectTagAdventure                   = "Adventure"
	ProjectTagFiction                     = "Fiction"
	ProjectTagNonfiction                  = "Nonfiction"
	ProjectTagHistorical                  = "Historical"
	ProjectTagDrama                       = "Drama"
	ProjectTagGames                       = "Games"
	ProjectTagStudy                       = "Study"
	ProjectTagHistory                     = "History"
	ProjectTagEnglish                     = "English"
	ProjectTagScience                     = "Science"
	ProjectTagMath                        = "Math"
	ProjectTagTraining                    = "Training"
)

func (a *ProjectTagArray) Scan(src interface{}) error {
	arr := common.StringArray{}
	err := arr.Scan(src)
	if err != nil {
		return err
	}
	newA := make(ProjectTagArray, len(arr.Val))
	for idx, v := range arr.Val {
		newA[idx] = ProjectTag(v)
	}
	*a = newA
	return nil
}

func (arr *ProjectTagArray) Value() (driver.Value, error) {
	v := []string{}
	for _, a := range *arr {
		v = append(v, fmt.Sprintf("\"%v\"", string(a)))
	}

	s := strings.Join(v, ",")
	return fmt.Sprintf("{%v}", s), nil
}

// Project is the model for a Workbench project
type Project struct {
	Model

	Title       string
	TeamID      uuid.UUID
	StartZoneID uuid.NullUUID // Expected Zone ID
	IsPrivate   bool
	Category    *ProjectCategory
	Tags        *ProjectTagArray

	Actors               []Actor                `db:"-"`
	Zones                []Zone                 `db:"-"`
	ZoneActors           []ZoneActor            `db:"-"`
	PrivateProjectGrants []PrivateProjectGrants `db:"-"`
	Notes                []Note                 `db:"-"`
}

type PublishedProject struct {
	ProjectID uuid.UUID
	TeamID    uuid.UUID
	CreatedAt gorp.NullTime `json:"CreatedAt,omitempty"`
}

func (p Project) PrepareMarshal() map[string]interface{} {
	result := map[string]interface{}{
		"ID":         p.ID,
		"Title":      p.Title,
		"CreatedAt":  p.CreatedAt.Time,
		"Zones":      p.Zones,
		"Actors":     p.Actors,
		"ZoneActors": p.ZoneActors,
		"Category":   p.Category,
		"Tags":       p.Tags,
	}

	if p.StartZoneID.Valid {
		result["StartZoneID"] = p.StartZoneID.UUID.String()
	}

	return result
}

type PrivateProjectGrants struct {
	ID        uuid.UUID `db:"ID, primarykey, autoincrement"`
	ProjectID uuid.UUID
	UserID    uuid.UUID
}

// AEID is an EntityID
// Useful for Redis key mapping
type AEID int

const (
	// AEIDActor EntityID for Actor
	AEIDActor AEID = iota

	// AEIDZone EntityID for Zone
	AEIDZone

	// AEIDTrigger EntityID for Trigger
	AEIDTrigger
	// AEIDDialogNode EntityID for DialogNode
	AEIDDialogNode
	// AEIDActionBundle EntityID for ActionBundle
	AEIDActionBundle
)

// DialogNode is a single instance of a Dialog
// If ParentNodes == nil, then it's the entry of a dialog
// If ChildNodes == nil, then it's the end of the dialog
type DialogNode struct {
	Model

	IsRoot    bool
	ProjectID uuid.UUID `json:"-"`
	ActorID   uuid.UUID `json:"-"`
	// Handles all other entry inputs
	UnknownHandler bool
	EntryInput     DialogInputArray
	RawLBlock
	ChildNodes  *[]*DialogNode `db:"-" json:"-"`
	ParentNodes *[]*DialogNode `db:"-" json:"-"`
}

func (a *DialogNode) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}

// DialogInput indicates valid dialog entry types
// As specified in https://talkative.ai
// Greeting (Example: “Hello <Actor>”)
// Provides an Actor
type DialogInput string

func (input DialogInput) Prepared() string {
	return string(input)
}

const (
	// For the "catch-all" unknown dialog handler
	DialogSpecialInputUnknown string = "[UNKNOWN]"
)

type DialogInputArray []DialogInput

func (a *DialogInputArray) Scan(src interface{}) error {
	arr := common.StringArray{}
	err := arr.Scan(src)
	if err != nil {
		return err
	}
	newA := make(DialogInputArray, len(arr.Val))
	for idx, v := range arr.Val {
		newA[idx] = DialogInput(v)
	}
	*a = newA
	return nil
}

func (arr *DialogInputArray) Value() (driver.Value, error) {
	v := []string{}
	for _, a := range *arr {
		v = append(v, fmt.Sprintf("\"%v\"", string(a)))
	}

	s := strings.Join(v, ",")
	return fmt.Sprintf("{%v}", s), nil
}

const (
	// DialogInputStatementVerb Verb statement
	// (Example: “I will <Verb> the <Actor>”)
	// Provides a Verb and an Actor
	DialogInputStatementVerb DialogInput = "statement_verb"
	// DialogInputGreeting Generic greeting
	// (Example: "Hello <Actor>")
	// Provides an optional Actor
	DialogInputGreeting DialogInput = "statement_greeting"
	// DialogInputFarewell Generic farewell
	// (Example: "Goodbye <Actor>")
	// Provides an optional Actor
	DialogInputFarewell DialogInput = "statement_farewell"
	// DialogInputQuestionVerb Verb question
	// (Example: “Did you <Verb> the <Actor>?”)
	// Provides a Verb and an Actor
	DialogInputQuestionVerb DialogInput = "question_verb"
	// DialogInputQuestionPossessional Possessional question
	// (Example: “Do you have <Actor>?”)
	// Provides an Actor
	DialogInputQuestionPossessional DialogInput = "question_possessional"
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

func (u *UUIDCreateID) UnmarshalJSON(text []byte) error {
	// Remove quotations
	var bytes []byte
	if len(text) > 1 && text[0] == '"' {
		bytes = text[1 : len(text)-1]
	}

	if strings.HasPrefix(string(bytes), "create") {
		str := string(bytes)
		u.CreateID = &str
		return nil
	}
	tmp := uuid.UUID{}
	err := tmp.UnmarshalText(bytes)
	if err != nil {
		return err
	}

	(*u).UUID = tmp
	return nil
}

// Actor model for the Actor entities
type Actor struct {
	Model

	Title           string
	ProjectID       uuid.UUID        `json:"-"`
	Dialogs         []DialogNode     `json:",omitempty" db:"-"`
	DialogRelations []DialogRelation `json:",omitempty" db:"-"`
}

type DialogRelation struct {
	ParentNodeID UUIDCreateID
	ChildNodeID  UUIDCreateID
	PatchAction  *PatchAction `json:",omitempty" db:"-"`
}

// Zone model for the Zone entities
type Zone struct {
	Model

	ProjectID   uuid.UUID `json:"-"`
	Title       string
	Description string
	Triggers    map[TriggerType]Trigger `db:"-"`
}

type ZoneActor struct {
	ZoneID      UUIDCreateID
	ActorID     UUIDCreateID
	PatchAction *PatchAction `json:",omitempty" db:"-"`
}

type TriggerType int

const (
	TriggerInitializeZone TriggerType = iota
	TriggerEnterZone
	TriggerExitZone
	TriggerVariableUpdate
)

// Trigger model for Trigger entities
type Trigger struct {
	TriggerType TriggerType
	ZoneID      UUIDCreateID
	RawLBlock
	PatchAction *PatchAction `json:",omitempty" db:"-"`
}

// Note model for the Note entities
type Note struct {
	Model
	Text string
}
