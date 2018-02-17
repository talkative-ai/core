package models

import (
	"net/url"
	"testing"
)

func TestActionSetIterable(t *testing.T) {
	AAS := ActionSet{}
	AAS.PlaySounds = make([]RAPlaySound, 2)
	AAS.PlaySounds[0].SoundType = RAPlaySoundTypeText
	AAS.PlaySounds[0].Val = "Hello world"
	AAS.PlaySounds[1].SoundType = RAPlaySoundTypeAudio
	AAS.PlaySounds[1].Val, _ = url.Parse("https://upload.wikimedia.org/wikipedia/commons/b/bb/Test_ogg_mp3_48kbps.wav")

	for a := range AAS.Iterable() {
		t.Log(a)
	}

}
