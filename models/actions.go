package models

import (
	"fmt"
	"net/url"

	"github.com/artificial-universe-maker/go-ssml"
)

type AumActionID uint32

const (
	SetGlobalVariables AumActionID = iota
	PlaySounds
	InitializeActorDialog
	SetZone
	ResetGame
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

	// Execute accepts two parameters.
	// The first parameter is the game state
	// The second parameter is the input parameters
	Execute(AumMutableRuntimeState, map[int32]interface{})
}

// ARAs should be passed both the existing game state, and the output SSML.
// Then, Execute() will mutate accordingly

type ARASetGlobalVariables map[int32]string

func (ara ARASetGlobalVariables) Compile() []byte {
	return []byte{}
}
func (ara ARASetGlobalVariables) Execute(state *AumMutableRuntimeState, params map[int32]interface{}) {
}

func (ara ARAPlaySounds) Compile() []byte {
	return []byte{}
}

type ARAPlaySoundsParams int32

const (
	ARAPlaySoundsParamText = iota
	ARAPlaySoundParamAudio
)

func (ara ARAPlaySounds) Execute(state *AumMutableRuntimeState, params map[int32]interface{}) {
	for k, v := range params {
		switch k {
		case ARAPlaySoundsParamText:
			state.OutputSSML.Text(v.(string))
			break
		case ARAPlaySoundParamAudio:
			state.OutputSSML.Audio(v.(*url.URL))
			break
		default:
			fmt.Println("Invalid ARAPlaySounds param", k)
			break
		}
	}
}

type ARAPlaySounds []int32
type ARAInitializeActorDialog int32
type ARASetZone int32
type ARAResetGame bool
