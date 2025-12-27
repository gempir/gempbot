package emotechief_test

import (
	"encoding/json"
	"testing"

	"github.com/gempir/gempbot/internal/channelpoint"
	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/emotechief"
	"github.com/gempir/gempbot/internal/emoteservice"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/store"
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
	helixClient := helixclient.NewMockClient()

	ec := emotechief.NewEmoteChief(config.NewMockConfig(), &store.Database{}, helixClient, emoteservice.NewSevenTvClient(store.NewMockStore()))

	opts := channelpoint.SevenTvAdditionalOptions{Slots: 1}
	marshalled, _ := json.Marshal(opts)

	redemption := helix.EventSubChannelPointsCustomRewardRedemptionEvent{
		UserInput: "bad input",
	}

	assert.False(t, ec.VerifySeventvRedemption(store.ChannelPointReward{AdditionalOptions: string(marshalled[:])}, redemption))
}

func TestCanVerifySevenTvEmoteRedemption(t *testing.T) {
	cfg := config.NewMockConfig()
	helixClient := helixclient.NewMockClient()

	ec := emotechief.NewEmoteChief(cfg, store.NewMockStore(), helixClient, emoteservice.NewMockApiClient())

	opts := channelpoint.SevenTvAdditionalOptions{Slots: 1}
	marshalled, _ := json.Marshal(opts)

	redemption := helix.EventSubChannelPointsCustomRewardRedemptionEvent{
		UserInput: "my emote :) https://7tv.app/emotes/60aed4fe423a803ccae373d3",
	}

	assert.True(t, ec.VerifySeventvRedemption(store.ChannelPointReward{AdditionalOptions: string(marshalled[:])}, redemption))
}

func TestCanHandleSevenTvEmoteRedemption(t *testing.T) {
	cfg := config.NewMockConfig()
	db := store.NewMockStore()
	helixClient := helixclient.NewClient(cfg, db)

	ec := emotechief.NewEmoteChief(cfg, db, helixClient, emoteservice.NewMockApiClient())

	opts := channelpoint.SevenTvAdditionalOptions{Slots: 1}
	marshalled, _ := json.Marshal(opts)

	redemption := helix.EventSubChannelPointsCustomRewardRedemptionEvent{
		UserInput: "my emote :) https://7tv.app/emotes/60aed4fe423a803ccae373d3",
	}

	ec.HandleSeventvRedemption(store.ChannelPointReward{AdditionalOptions: string(marshalled[:])}, redemption, true)
}
