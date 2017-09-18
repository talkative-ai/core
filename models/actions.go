package models

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"log"

	"github.com/artificial-universe-maker/go-ssml"
)

// AumActionID is an ID for each "action" type
// Used in Lakshmi compilation process
// And in Brahman runtime process
type AumActionID uint64

const (
	// AAIDSetGlobalVariable AumActionID for SetGlobalVariable
	AAIDSetARVariable AumActionID = iota
	// AAIDPlaySound AumActionID for PlaySound
	AAIDPlaySound
	// AAIDInitializeActorDialog AumActionID for InitializeActorDialog
	AAIDInitializeActorDialog
	// AAIDSetZone AumActionID for SetZone
	AAIDSetZone
	// AAIDResetGame AumActionID for ResetGame
	AAIDResetGame
)

// AumActionSet is a pre-bundled set of actions
// These actions either mutate the runtime state or mutate the output dialog
type AumActionSet struct {
	SetGlobalVariables    []ARASetVariable
	PlaySounds            []ARAPlaySound
	InitializeActorDialog int32
	SetZone               ARASetZone
	ResetGame             bool
}

func (a *AumActionSet) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &a)
}

func (a *AumActionSet) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Iterable will output all of the AumRuntimeActions within the AumActionSet
// This is useful for easily bundling actions within Lakshmi
// Without having to create ad hoc functions
func (AAS AumActionSet) Iterable() <-chan AumRuntimeAction {
	ch := make(chan AumRuntimeAction)
	go func() {
		defer close(ch)
		// TODO: Add other actions in here
		for _, r := range AAS.PlaySounds {
			action := r
			ch <- &action
		}
		for _, r := range AAS.SetGlobalVariables {
			action := r
			ch <- &action
		}

		if uint64(AAS.SetZone) != 0 {
			ch <- &AAS.SetZone
		}
	}()
	return ch
}

// AumMutableRuntimeState is used by Brahman
// It contains the current State of the running game / project (called the runtime state)
// and the OutputSSML, Speech-synthesis markup language
type AumMutableRuntimeState struct {
	State      MutableRuntimeState
	OutputSSML ssml.Builder
}

type MutableRuntimeState struct {
	Zone          string
	PubID         string
	CurrentDialog *string
	ZoneActors    map[string][]string
	ARVariables   map[string]*ARVariable
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

func (a *MutableRuntimeState) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// AumRuntimeAction is an interface for all the actions within an AumActionSet
// Combined with the AumActionSet Iterable(), compilation is easy
type AumRuntimeAction interface {
	// Compile is used by Lakshmi
	// Returns the compiled []byte slice of the runtime action
	Compile() []byte

	// CreateFrom is useful for testing
	CreateFrom([]byte) error

	// GetAAID returns the AumActionID of the current RuntimeAction
	// Useful for Lakshmi
	GetAAID() AumActionID

	// Execute will mutate the AumMutableRuntimeState in some way
	// Whether it's the state itself or the OutputSSML
	Execute(*AumMutableRuntimeState)
}

// ARAPlaySoundType is an enum of different ARAPlaySound types
// (Speech-synthesis or Audio File)
type ARAPlaySoundType uint8

const (
	// ARAPlaySoundTypeText Speech-synthesis
	ARAPlaySoundTypeText ARAPlaySoundType = iota
	// ARAPlaySoundTypeAudio URL to an audio file
	ARAPlaySoundTypeAudio
)

// ARAPlaySound AumRuntimeAction PlaySound
// This action mutates the OutputSSML of the AumMutableRuntimeState
type ARAPlaySound struct {
	SoundType ARAPlaySoundType
	Val       interface{}
}

func GetActionFromID(id AumActionID) AumRuntimeAction {
	switch id {
	case AAIDPlaySound:
		return &ARAPlaySound{}
	case AAIDSetZone:
		n := ARASetZone(0)
		return &n
	case AAIDSetARVariable:
		return &ARASetVariable{}
	default:
		return nil
	}
}

//////////////////
// ARAPlaySound //
//////////////////

// GetAAID returns the AumActionID of the current RuntimeAction
func (ara ARAPlaySound) GetAAID() AumActionID {
	return AAIDPlaySound
}

// Compile is used by Lakshmi
// Returns the compiled []byte slice of the runtime action
// To be stored in Redis
func (ara ARAPlaySound) Compile() []byte {
	compiled := []byte{}
	compiled = append(compiled, byte(ara.SoundType))
	switch ara.SoundType {
	case ARAPlaySoundTypeText:
		compiled = append(compiled, []byte(ara.Val.(string))...)
		break
	case ARAPlaySoundTypeAudio:
		compiled = append(compiled, []byte(ara.Val.(*url.URL).String())...)
		break
	}

	return compiled
}

// CreateFrom is used for evaluating the actions in Brahman and followed by Execute
// This could be put in a single "Execute" but this is less monolothic
func (ara *ARAPlaySound) CreateFrom(bytes []byte) error {
	ara.SoundType = ARAPlaySoundType(bytes[0])
	bytes = bytes[1:]
	switch ara.SoundType {
	case ARAPlaySoundTypeText:
		ara.Val = string(bytes)
		break
	case ARAPlaySoundTypeAudio:
		var err error
		ara.Val, err = url.Parse(string(bytes))
		if err != nil {
			return err
		}
	}
	return nil
}

// Execute will mutate the AumMutableRuntimeState in some way
// Whether it's the state itself or the OutputSSML
func (ara ARAPlaySound) Execute(state *AumMutableRuntimeState) {
	switch ara.SoundType {
	case ARAPlaySoundTypeText:
		r, err := regexp.Compile(`{{\w+}}`)
		v := ara.Val.(string)
		if err != nil {
			log.Println("[ERROR] Inavalid type on ARASetVariable Execute")
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
	case ARAPlaySoundTypeAudio:
		state.OutputSSML = state.OutputSSML.Audio(ara.Val.(*url.URL))
		break
	}
}

////////////////
// ARASetZone //
////////////////
type ARASetZone uint64

// GetAAID returns the AumActionID of the current RuntimeAction
func (ara *ARASetZone) GetAAID() AumActionID {
	return AAIDSetZone
}

// Compile is used by Lakshmi
// Returns the compiled []byte slice of the runtime action
// To be stored in Redis
func (ara ARASetZone) Compile() []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(ara))
	return b
}

// CreateFrom is used for evaluating the actions in Brahman and followed by Execute
// This could be put in a single "Execute" but this is less monolothic
func (ara *ARASetZone) CreateFrom(bytes []byte) error {
	n := ARASetZone(binary.LittleEndian.Uint64(bytes))
	*ara = n
	return nil
}

// Execute will mutate the AumMutableRuntimeState in some way
// Whether it's the state itself or the OutputSSML
func (ara *ARASetZone) Execute(state *AumMutableRuntimeState) {
	state.State.Zone = fmt.Sprintf("%v", *ara)
	state.State.CurrentDialog = nil
}

////////////////////
// ARASetVariable //
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

type ARASetVariable struct {
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

// GetAAID returns the AumActionID of the current RuntimeAction
func (ara *ARASetVariable) GetAAID() AumActionID {
	return AAIDSetARVariable
}

// Compile is used by Lakshmi
// Returns the compiled []byte slice of the runtime action
// To be stored in Redis
func (ara ARASetVariable) Compile() []byte {
	b := []byte{}
	return b
}

// CreateFrom is used for evaluating the actions in Brahman and followed by Execute
// This could be put in a single "Execute" but this is less monolothic
func (ara *ARASetVariable) CreateFrom(bytes []byte) error {
	return nil
}

// Execute will mutate the AumMutableRuntimeState in some way
// Whether it's the state itself or the OutputSSML
func (ara *ARASetVariable) Execute(state *AumMutableRuntimeState) {
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
			log.Println("[ERROR] Inavalid type on ARASetVariable Execute")
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
			log.Println("[ERROR] Inavalid type on ARASetVariable Execute")
			return
		}
	case SVOSubtract:
		switch o := original.(type) {
		case int64:
			newval = o - n.(int64)
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Println("[ERROR] Inavalid type on ARASetVariable Execute")
			return
		}
	case SVODivide:
		switch o := original.(type) {
		case int64:
			newval = o / n.(int64)
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Println("[ERROR] Inavalid type on ARASetVariable Execute")
			return
		}
	case SVOModulo:
		switch o := original.(type) {
		case int64:
			newval = o % n.(int64)
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Println("[ERROR] Inavalid type on ARASetVariable Execute")
			return
		}
	case SVONot:
		switch o := original.(type) {
		case bool:
			newval = !o
		default:
			// TODO: Better error handling here
			log.SetFlags(log.Llongfile | log.Ltime)
			log.Println("[ERROR] Inavalid type on ARASetVariable Execute")
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
			log.Println("[ERROR] Inavalid type on ARASetVariable Execute")
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
			log.Println("[ERROR] Inavalid type on ARASetVariable Execute")
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
			log.Println("[ERROR] Inavalid type on ARASetVariable Execute")
			return
		}
	}
	state.State.ARVariables[ara.Target].Val = newval
}
