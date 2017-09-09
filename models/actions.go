package models

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"net/url"

	"github.com/artificial-universe-maker/go-ssml"
)

// AumActionID is an ID for each "action" type
// Used in Lakshmi compilation process
// And in Brahman runtime process
type AumActionID uint64

const (
	// AAIDSetGlobalVariable AumActionID for SetGlobalVariable
	AAIDSetGlobalVariable AumActionID = iota
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
	SetGlobalVariables    map[int32]string
	PlaySounds            []ARAPlaySound
	InitializeActorDialog int32
	SetZone               int32
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

// GetAAID returns the AumActionID of the current RuntimeAction
func (ara ARAPlaySound) GetAAID() AumActionID {
	return AAIDPlaySound
}

func GetActionFromID(id AumActionID) AumRuntimeAction {
	switch id {
	case AAIDPlaySound:
		return &ARAPlaySound{}
	default:
		return nil
	}
}

// Compile is used by Lakshmi
// Returns the compiled []byte slice of the runtime action
func (ara ARAPlaySound) Compile() []byte {
	compiled := []byte{}
	b := make([]byte, 4)
	compiled = append(compiled, byte(ara.SoundType))
	switch ara.SoundType {
	case ARAPlaySoundTypeText:
		compiled = append(compiled, []byte(ara.Val.(string))...)
		break
	case ARAPlaySoundTypeAudio:
		compiled = append(compiled, []byte(ara.Val.(*url.URL).String())...)
		break
	}

	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(compiled)))
	finished := []byte{}
	finished = append(finished, b...)
	finished = append(finished, compiled...)

	return finished
}

// Execute will mutate the AumMutableRuntimeState in some way
// Whether it's the state itself or the OutputSSML
func (ara ARAPlaySound) Execute(state *AumMutableRuntimeState) {
	switch ara.SoundType {
	case ARAPlaySoundTypeText:
		state.OutputSSML = state.OutputSSML.Paragraph(ara.Val.(string))
		break
	case ARAPlaySoundTypeAudio:
		state.OutputSSML = state.OutputSSML.Audio(ara.Val.(*url.URL))
		break
	}
}

// CreateFrom is used for evaluating the actions in Brahman and followed by Execute
// This could be put in a single "Execute" but this way is less monolothic
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
