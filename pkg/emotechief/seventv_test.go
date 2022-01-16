package emotechief

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanGetSevenTvEmoteFromMessage(t *testing.T) {
	tests := []struct {
		message string
		emoteId string
	}{
		{"some message https://7tv.app/emotes/60ccf4479f5edeff9938fa77 some more message", "60ccf4479f5edeff9938fa77"},
		{"https://7tv.app/emotes/60aed4fe423a803ccae373d3", "60aed4fe423a803ccae373d3"},
		{"some message", ""},
	}

	for _, test := range tests {
		emote, err := GetSevenTvEmoteId(test.message)
		if err != nil && test.emoteId != "" {
			t.Error(err.Error())
		}

		assert.Equal(t, test.emoteId, emote, "could not parse emoteId")
	}

}
