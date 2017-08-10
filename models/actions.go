package models

import (
	"encoding/binary"
	"net/url"

	"github.com/artificial-universe-maker/go-ssml"
)

type AumActionID uint64

const (
	AAIDSetGlobalVariable AumActionID = iota
	AAIDPlaySound
	AAIDInitializeActorDialog
	AAIDSetZone
	AAIDResetGame
)

type AumActionSet struct {
	SetGlobalVariables    map[int32]string
	PlaySounds            []int32
	InitializeActorDialog int32
	SetZone               int32
	ResetGame             bool
}

type AumMutableRuntimeState struct {
	State      map[string]string
	OutputSSML ssml.Builder
}

type AumRuntimeAction interface {
	Compile() []byte
	CreateFrom([]byte)

	// Execute accepts two parameters.
	// The first parameter is the game state
	// The second parameter is the input parameters
	Execute(AumMutableRuntimeState)
}

type ARAPlaySoundType uint8

const (
	ARAPlaySoundTypeText ARAPlaySoundType = iota
	ARAPlaySoundTypeAudio
)

type ARAPlaySound struct {
	SoundType ARAPlaySoundType
	Value     interface{}
}
type ARAInitializeActorDialog int32
type ARASetZone int32
type ARAResetGame bool

// ARAs should be passed both the existing game state, and the output SSML.
// Then, Execute() will mutate accordingly

func (ara ARAPlaySound) Compile() []byte {
	compiled := []byte{}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(AAIDPlaySound))
	compiled = append(compiled, b...)
	compiled = append(compiled, byte(ara.SoundType))
	switch ara.SoundType {
	case ARAPlaySoundTypeText:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(len(ara.Value.(string))))
		compiled = append(compiled, b...)
		compiled = append(compiled, []byte(ara.Value.(string))...)
		break
	case ARAPlaySoundTypeAudio:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(len(ara.Value.(*url.URL).String())))
		compiled = append(compiled, b...)
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

func (ara ARAPlaySound) CreateFrom(bytes []byte) {

}
