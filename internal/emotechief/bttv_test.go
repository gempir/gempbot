package emotechief_test

import (
	"testing"

	"github.com/gempir/gempbot/internal/emotechief"
	"github.com/stretchr/testify/assert"
)

func TestCanGetBttvTvEmoteFromMessage(t *testing.T) {
	tests := []struct {
		message string
		emoteId string
	}{
		{"some message https://betterttv.com/emotes/59f27b3f4ebd8047f54dee29 some more message", "59f27b3f4ebd8047f54dee29"},
		{"https://betterttv.com/emotes/5d20a55de1cfde376e532972", "5d20a55de1cfde376e532972"},
		{"some message", ""},
	}

	for _, test := range tests {
		emote, err := emotechief.GetBttvEmoteId(test.message)
		if err != nil && test.emoteId != "" {
			t.Error(err.Error())
		}

		assert.Equal(t, test.emoteId, emote, "could not parse emoteId")
	}

}
