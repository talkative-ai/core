package models

import (
	"fmt"
	"net/url"
	"testing"

	ssml "github.com/talkative-ai/go-ssml"
	"github.com/talkative-ai/lakshmi/prepare"
)

func TestActionBundleEval(t *testing.T) {
	AAS := ActionSet{}
	AAS.PlaySounds = make([]RAPlaySound, 2)
	AAS.PlaySounds[0].SoundType = RAPlaySoundTypeText
	AAS.PlaySounds[0].Val = "Hello world"
	AAS.PlaySounds[1].SoundType = RAPlaySoundTypeAudio
	AAS.PlaySounds[1].Val, _ = url.Parse("https://upload.wikimedia.org/wikipedia/commons/b/bb/Test_ogg_mp3_48kbps.wav")

	runtimeState := AIRequest{
		State:      MutableAIRequestState{},
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
