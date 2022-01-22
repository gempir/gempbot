package emotechief_test

import (
	"encoding/json"
	"testing"

	"github.com/gempir/gempbot/pkg/channelpoint"
	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/emotechief"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/nicklaw5/helix/v2"
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
		emote, err := emotechief.GetSevenTvEmoteId(test.message)
		if err != nil && test.emoteId != "" {
			t.Error(err.Error())
		}

		assert.Equal(t, test.emoteId, emote, "could not parse emoteId")
	}

}

func TestCanNotVerifySevenTvEmoteRedemption(t *testing.T) {
	ec := emotechief.NewEmoteChief(config.NewTestConfig(), &store.Database{}, &helixclient.Client{}, chat.NewClient(config.NewTestConfig()))

	opts := channelpoint.BttvAdditionalOptions{Slots: 1}
	marshalled, _ := json.Marshal(opts)

	redemption := helix.EventSubChannelPointsCustomRewardRedemptionEvent{
		UserInput: "bad input",
	}

	assert.False(t, ec.VerifySeventvRedemption(store.ChannelPointReward{AdditionalOptions: string(marshalled[:])}, redemption))
}
