package emotechief

import "github.com/gempir/bitraft/pkg/log"

func (e *EmoteChief) SetSevenTvEmote(channelUserID, emoteId, channel string, slots int) (addedEmote *bttvEmoteResponse, removedEmote *bttvEmoteResponse, err error) {

	log.Info("adding emote ", emoteId)

	return nil, nil, nil
}
