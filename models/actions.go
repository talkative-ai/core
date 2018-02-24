package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/talkative-ai/go.uuid"

	"log"

	"github.com/talkative-ai/core/redis"
	"github.com/talkative-ai/go-ssml"
)

// ActionID is an ID for each "action" type
// Used in Lakshmi compilation process
// And in Brahman runtime process
type ActionID uint64

const (
	// RAIDSetGlobalVariable ActionID for SetGlobalVariable
	RAIDSetARVariable ActionID = iota
	// RAIDPlaySound ActionID for PlaySound
	RAIDPlaySound
	// RAIDInitializeActorDialog ActionID for InitializeActorDialog
	RAIDInitializeActorDialog
	// RAIDSetZone ActionID for SetZone
	RAIDSetZone
	// RAIDResetApp ActionID for ResetApp
	RAIDResetApp
)

// ActionSet is a pre-bundled set of actions
// These actions either mutate the runtime state or mutate the output dialog
type ActionSet struct {
	SetGlobalVariables    []RASetVariable
	PlaySounds            []RAPlaySound
	InitializeActorDialog uuid.UUID
	SetZone               RASetZone
	ResetApp              RAResetApp
}

func (a *ActionSet) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}

func (a *ActionSet) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Iterable will output all of the RequestActions within the ActionSet
// This is useful for easily bundling actions within Lakshmi
// Without having to create ad hoc functions
func (AAS ActionSet) Iterable() <-chan RequestAction {
	ch := make(chan RequestAction)
	go func() {
		defer close(ch)
		for _, r := range AAS.PlaySounds {
			action := r
			ch <- &action
		}

		for _, r := range AAS.SetGlobalVariables {
			action := r
			ch <- &action
		}

		if uuid.UUID(AAS.SetZone) != uuid.Nil {
			ch <- &AAS.SetZone
		}

		if AAS.ResetApp {
			ch <- &AAS.ResetApp
		}
	}()
	return ch
}

// AIRequest is used by Brahman
// It contains the current State of the running game / project (called the runtime state)
// and the OutputSSML, Speech-synthesis markup language
type AIRequest struct {
	State      MutableAIRequestState
	OutputSSML ssml.Builder
}

type MutableAIRequestState struct {
	Zone            uuid.UUID
	ProjectID       uuid.UUID
	PubID           string
	CurrentDialog   *string
	ZoneActors      map[uuid.UUID][]string
	ZoneInitialized map[uuid.UUID]bool
	ARVariables     map[string]*ARVariable
	Demo            bool
}

type ARVariable struct {
	T   string
	Val interface{}
}

func (arv *ARVariable) Get() interface{} {
	var v interface{}
	switch arv.T {
	case "int":
		v = arv.Val.(int)
	case "bool":
		v = arv.Val.(bool)
	case "array":
		v = arv.Val.([]ARVariable)
	default:
		v = fmt.Sprintf("%+v", arv.Val)
	}
	return v
}

func (a *MutableAIRequestState) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// RequestAction is an interface for all the actions within an ActionSet
// Combined with the ActionSet Iterable(), compilation is easy
type RequestAction interface {
	// Compile is used by Lakshmi
	// Returns the compiled []byte slice of the runtime action
	Compile() []byte

	// CreateFrom is useful for testing
	CreateFrom([]byte) error

	// GetRAID returns the ActionID of the current RequestAction
	// Useful for Lakshmi
	GetRAID() ActionID

	// Execute will mutate the AIRequest in some way
	// Whether it's the state itself or the OutputSSML
	Execute(*AIRequest)
}

// RAPlaySoundType is an enum of different RAPlaySound types
// (Speech-synthesis or Audio File)
type RAPlaySoundType uint8

const (
	// RAPlaySoundTypeText Speech-synthesis
	RAPlaySoundTypeText RAPlaySoundType = iota
	// RAPlaySoundTypeAudio URL to an audio file
	RAPlaySoundTypeAudio
)

// RAPlaySound RequestAction PlaySound
// This action mutates the OutputSSML of the AIRequest
type RAPlaySound struct {
	SoundType RAPlaySoundType
	Val       interface{}
}

func GetActionFromID(id ActionID) RequestAction {
	switch id {
	case RAIDPlaySound:
		return &RAPlaySound{}
	case RAIDSetZone:
		n := RASetZone(uuid.Nil)
		return &n
	case RAIDSetARVariable:
		return &RASetVariable{}
	case RAIDResetApp:
		n := RAResetApp(false)
		return &n
	default:
		log.Fatalln("Unsupported action id:", id)
		return nil
	}
}

//////////////////
// RAPlaySound //
//////////////////

// GetRAID returns the ActionID of the current RequestAction
func (ara RAPlaySound) GetRAID() ActionID {
	return RAIDPlaySound
}

// Compile is used by Lakshmi
// Returns the compiled []byte slice of the runtime action
// To be stored in Redis
func (ara RAPlaySound) Compile() []byte {
	compiled := []byte{}
	compiled = append(compiled, byte(ara.SoundType))
	switch ara.SoundType {
	case RAPlaySoundTypeText:
		compiled = append(compiled, []byte(ara.Val.(string))...)
		break
	case RAPlaySoundTypeAudio:
		compiled = append(compiled, []byte(ara.Val.(*url.URL).String())...)
		break
	}

	return compiled
}

// CreateFrom is used for evaluating the actions in Brahman and followed by Execute
// This could be put in a single "Execute" but this is less monolothic
func (ara *RAPlaySound) CreateFrom(bytes []byte) error {
	ara.SoundType = RAPlaySoundType(bytes[0])
	bytes = bytes[1:]
	switch ara.SoundType {
	case RAPlaySoundTypeText:
		ara.Val = string(bytes)
		break
	case RAPlaySoundTypeAudio:
		var err error
		ara.Val, err = url.Parse(string(bytes))
		if err != nil {
			return err
		}
	}
	return nil
}

// Execute will mutate the AIRequest in some way
// Whether it's the state itself or the OutputSSML
func (ara RAPlaySound) Execute(state *AIRequest) {
	switch ara.SoundType {
	case RAPlaySoundTypeText:
		r, err := regexp.Compile(`{{\w+}}`)
		v := ara.Val.(string)
		if err != nil {
			log.Fatal("[ERROR] Inavalid type on RASetVariable Execute")
			return
		}
		newv := r.ReplaceAllFunc([]byte(v), func(b []byte) []byte {
			v := b[2 : len(b)-2]
			// TODO: Support array indices
			val := state.State.ARVariables[string(v)].Val
			return []byte(fmt.Sprintf("%v", val))
		})
		state.OutputSSML = state.OutputSSML.Paragraph(string(newv))
		break
	case RAPlaySoundTypeAudio:
		state.OutputSSML = state.OutputSSML.Audio(ara.Val.(*url.URL))
		break
	}
}

////////////////
// RASetZone //
////////////////
type RASetZone uuid.UUID

// GetRAID returns the ActionID of the current RequestAction
func (ara *RASetZone) GetRAID() ActionID {
	return RAIDSetZone
}

// Compile is used by Lakshmi
// Returns the compiled []byte slice of the runtime action
// To be stored in Redis
func (ara RASetZone) Compile() []byte {
	return uuid.UUID(ara).Bytes()
}

// CreateFrom is used for evaluating the actions in Brahman and followed by Execute
// This could be put in a single "Execute" but this is less monolothic
func (ara *RASetZone) CreateFrom(bytes []byte) error {
	*ara = RASetZone(uuid.FromBytesOrNil(bytes))
	return nil
}

func (ara *RASetZone) String() string {
	return uuid.UUID(*ara).String()
}

func (ara *RASetZone) UUID() uuid.UUID {
	return uuid.UUID(*ara)
}

func (ara RASetZone) MarshalText() (text []byte, err error) {
	text = []byte(ara.String())
	return
}

func (ara *RASetZone) UnmarshalText(text []byte) error {
	tmp := uuid.UUID{}
	err := tmp.UnmarshalText(text)
	if err != nil {
		return err
	}

	*ara = RASetZone(tmp)
	return nil
}

// Execute will mutate the AIRequest in some way
// Whether it's the state itself or the OutputSSML
func (ara *RASetZone) Execute(message *AIRequest) {
	message.State.Zone = ara.UUID()
	message.State.CurrentDialog = nil

	if message.State.ZoneInitialized[message.State.Zone] {
		return
	}

	message.State.ZoneInitialized[message.State.Zone] = true

	res := redis.Instance.HGet(
		KeynavCompiledTriggersWithinZone(message.State.PubID, ara.String()),
		fmt.Sprintf("%v", TriggerInitializeZone)).Val()

	// There is no initialize trigger
	if res == "" {
		return
	}

	stateComms := make(chan AIRequest)
	result := LogicLazyEval(stateComms, []byte(res))
	for res := range result {
		if res.Error != nil {
			log.Fatal("Error in SetZone with logic evaluation", res.Error)
			return
		}
		bundleBinary, err := redis.Instance.Get(res.Value).Bytes()
		if err != nil {
			log.Fatal("Error in SetZone fetching action bundle binary", err)
			return
		}
		err = ActionBundleEval(message, bundleBinary)
		if err != nil {
			log.Fatal("Error in SetZone processing action bundle binary", err)
			return
		}
	}
}

////////////////
// RASetZone //
////////////////
type RAResetApp bool

// GetRAID returns the ActionID of the current RequestAction
func (ara *RAResetApp) GetRAID() ActionID {
	return RAIDResetApp
}

// Compile is used by Lakshmi
// Returns the compiled []byte slice of the runtime action
// To be stored in Redis
func (ara RAResetApp) Compile() []byte {
	return []byte{}
}

// CreateFrom is used for evaluating the actions in Brahman and followed by Execute
// This could be put in a single "Execute" but this is less monolothic
func (ara *RAResetApp) CreateFrom(bytes []byte) error {
	*ara = true
	return nil
}

// Execute will mutate the AIRequest in some way
// Whether it's the state itself or the OutputSSML
func (ara *RAResetApp) Execute(message *AIRequest) {
	if *ara {
		// The reset is happening from inside the app
	} else {
		// The reset is being triggered manually
	}
	message.State.ZoneActors = map[uuid.UUID][]string{}
	message.State.ZoneInitialized = map[uuid.UUID]bool{}
	for _, zoneID := range redis.Instance.SMembers(
		fmt.Sprintf("%v:%v", KeynavProjectMetadataStatic(message.State.PubID), "all_zones")).Val() {
		zUUID := uuid.FromStringOrNil(zoneID)
		message.State.ZoneActors[zUUID] =
			redis.Instance.SMembers(KeynavCompiledActorsWithinZone(message.State.PubID, zoneID)).Val()
		message.State.ZoneInitialized[zUUID] = false
	}
	zoneID := redis.Instance.HGet(KeynavProjectMetadataStatic(message.State.PubID), "start_zone_id").Val()
	setZone := RASetZone(uuid.FromStringOrNil(zoneID))
	setZone.Execute(message)
}

////////////////////
// RASetVariable //
////////////////////
type SetVariableOperation int

const (
	// SVOSet is for any
	SVOSet SetVariableOperation = iota

	// SVOAdd is for int/string
	SVOAdd

	// SVOSubtract is for int
	SVOSubtract
	// SVODivide is for int
	SVODivide
	// SVOMultiply is for int
	SVOMultiply
	// SVOModulo is for int
	SVOModulo

	// SVONot is for bool
	SVONot

	// SVOInsert is for array
	SVOInsert
	// SVODelete is for array
	SVODelete

	// SVOReplace is for string
	SVOReplace
)

type RASetVariable struct {
	Target    string
	Operation SetVariableOperation
	With      ParametizedARVariable
}

type ParametizedARVariable struct {
	// If Key == nil then the variable is inlined
	Key        *string
	Params     map[string]interface{}
	ARVariable *ARVariable
}

// GetRAID returns the ActionID of the current RequestAction
func (ara *RASetVariable) GetRAID() ActionID {
	return RAIDSetARVariable
}

// Compile is used by Lakshmi
// Returns the compiled []byte slice of the runtime action
// To be stored in Redis
func (ara RASetVariable) Compile() []byte {
	b := []byte{}
	return b
}

// CreateFrom is used for evaluating the actions in Brahman and followed by Execute
// This could be put in a single "Execute" but this is less monolothic
func (ara *RASetVariable) CreateFrom(bytes []byte) error {
	return nil
}

// Execute will mutate the AIRequest in some way
// Whether it's the state itself or the OutputSSML
func (ara *RASetVariable) Execute(state *AIRequest) {
	original := state.State.ARVariables[ara.Target].Get()
	var n interface{}
	if ara.With.Key != nil {
		n = state.State.ARVariables[*ara.With.Key].Get()
	} else {
		n = ara.With.ARVariable.Get()
	}
	var newval interface{}

	switch ara.Operation {
	case SVOSet:
		if ara.With.ARVariable.T != state.State.ARVariables[ara.Target].T {
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Fatal("[ERROR] Inavalid type on RASetVariable Execute")
			return
		}
		newval = n
	case SVOAdd:
		switch o := original.(type) {
		case int64:
			newval = o + n.(int64)
		case string:
			newval = fmt.Sprintf("%v%v", o, n.(string))
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Fatal("[ERROR] Inavalid type on RASetVariable Execute")
			return
		}
	case SVOSubtract:
		switch o := original.(type) {
		case int64:
			newval = o - n.(int64)
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Fatal("[ERROR] Inavalid type on RASetVariable Execute")
			return
		}
	case SVODivide:
		switch o := original.(type) {
		case int64:
			newval = o / n.(int64)
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Fatal("[ERROR] Inavalid type on RASetVariable Execute")
			return
		}
	case SVOModulo:
		switch o := original.(type) {
		case int64:
			newval = o % n.(int64)
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Fatal("[ERROR] Inavalid type on RASetVariable Execute")
			return
		}
	case SVONot:
		switch o := original.(type) {
		case bool:
			newval = !o
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Fatal("[ERROR] Inavalid type on RASetVariable Execute")
			return
		}
	case SVOInsert:
		index := ara.With.Params["Index"].(int)
		switch o := original.(type) {
		case []ARVariable:
			newval = []ARVariable{}
			newval = append(newval.([]ARVariable), o[:index]...)
			newval = append(newval.([]ARVariable), n.(ARVariable))
			newval = append(newval.([]ARVariable), o[index:]...)
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Fatal("[ERROR] Inavalid type on RASetVariable Execute")
			return
		}
	case SVODelete:
		index := ara.With.Params["Index"].(int)
		switch o := original.(type) {
		case []ARVariable:
			newval = []ARVariable{}
			newval = append(o[:index], o[:index+1]...)
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Fatal("[ERROR] Inavalid type on RASetVariable Execute")
			return
		}
	case SVOReplace:
		search := ara.With.Params["Search"].(string)
		replace := ara.With.Params["Replace"].(string)
		switch o := original.(type) {
		case string:
			newval = strings.Replace(o, search, replace, 0)
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Fatal("[ERROR] Inavalid type on RASetVariable Execute")
			return
		}
	}
	state.State.ARVariables[ara.Target].Val = newval
}
