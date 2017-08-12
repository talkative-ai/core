package models

import (
	"encoding/binary"
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

// Iterable will output all of the AumRuntimeActions within the AumActionSet
// This is useful for easily bundling actions within Lakshmi
// Without having to create ad hoc functions
func (AAS AumActionSet) Iterable() <-chan AumRuntimeAction {
	ch := make(chan AumRuntimeAction)
	go func() {
		defer func() {
			close(ch)
		}()
		for _, r := range AAS.PlaySounds {
			ch <- &r
		}
	}()
	return ch
}

// AumMutableRuntimeState is used by Brahman
// It contains the current State of the running game / project (called the runtime state)
// and the OutputSSML, Speech-synthesis markup language
type AumMutableRuntimeState struct {
	State      map[string]string
	OutputSSML ssml.Builder
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
	Value     interface{}
}

// GetAAID returns the AumActionID of the current RuntimeAction
func (ara ARAPlaySound) GetAAID() AumActionID {
	return AAIDPlaySound
}

// Compile is used by Lakshmi
// Returns the compiled []byte slice of the runtime action
func (ara ARAPlaySound) Compile() []byte {
	compiled := []byte{}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(AAIDPlaySound))
	compiled = append(compiled, b...)
	compiled = append(compiled, byte(ara.SoundType))
	switch ara.SoundType {
	case ARAPlaySoundTypeText:
		compiled = append(compiled, []byte(ara.Value.(string))...)
		break
	case ARAPlaySoundTypeAudio:
		compiled = append(compiled, []byte(ara.Value.(*url.URL).String())...)
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
		state.OutputSSML = state.OutputSSML.Text(ara.Value.(string))
		break
	case ARAPlaySoundTypeAudio:
		state.OutputSSML = state.OutputSSML.Audio(ara.Value.(*url.URL))
		break
	}
}

// CreateFrom is useful for testing
func (ara *ARAPlaySound) CreateFrom(bytes []byte) error {
	ara.SoundType = ARAPlaySoundType(bytes[0])
	bytes = bytes[1:]
	switch ara.SoundType {
	case ARAPlaySoundTypeText:
		ara.Value = string(bytes[:])
		break
	case ARAPlaySoundTypeAudio:
		var err error
		ara.Value, err = url.Parse(string(bytes[:]))
		if err != nil {
			return err
		}
	}
	return nil
}
