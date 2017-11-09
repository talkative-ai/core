package models

import (
	"net/url"
	"testing"
)

func TestAumActionSetIterable(t *testing.T) {
	AAS := AumActionSet{}
	AAS.PlaySounds = make([]ARAPlaySound, 2)
	AAS.PlaySounds[0].SoundType = ARAPlaySoundTypeText
	AAS.PlaySounds[0].Val = "Hello world"
	AAS.PlaySounds[1].SoundType = ARAPlaySoundTypeAudio
	AAS.PlaySounds[1].Val, _ = url.Parse("https://upload.wikimedia.org/wikipedia/commons/b/bb/Test_ogg_mp3_48kbps.wav")

	for a := range AAS.Iterable() {
		t.Log(a)
	}

}
