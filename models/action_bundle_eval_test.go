package models

import (
	"fmt"
	"net/url"
	"testing"

	ssml "github.com/artificial-universe-maker/go-ssml"
	"github.com/artificial-universe-maker/lakshmi/prepare"
)

func TestActionBundleEval(t *testing.T) {
	AAS := AumActionSet{}
	AAS.PlaySounds = make([]ARAPlaySound, 2)
	AAS.PlaySounds[0].SoundType = ARAPlaySoundTypeText
	AAS.PlaySounds[0].Val = "Hello world"
	AAS.PlaySounds[1].SoundType = ARAPlaySoundTypeAudio
	AAS.PlaySounds[1].Val, _ = url.Parse("https://upload.wikimedia.org/wikipedia/commons/b/bb/Test_ogg_mp3_48kbps.wav")

	runtimeState := AumMutableRuntimeState{
		State:      MutableRuntimeState{},
		OutputSSML: ssml.NewBuilder(),
	}

	bundled := prepare.BundleActions(AAS)

	err := ActionBundleEval(&runtimeState, bundled)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if runtimeState.OutputSSML.String() != `<speak>Hello world<audio src="https://upload.wikimedia.org/wikipedia/commons/b/bb/Test_ogg_mp3_48kbps.wav" /></speak>` {
		fmt.Println("Unexpected runtimeState Output SSML")
		t.Fail()
	}
}
