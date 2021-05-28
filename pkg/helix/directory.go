package helix

import (
	"github.com/gempir/bitraft/pkg/log"
	"github.com/nicklaw5/helix"
)

func (c *Client) GetTopChannels() []string {
	response, err := c.Client.GetStreams(&helix.StreamsParams{Type: "live", First: 100})
	if err != nil {
		log.Error(err.Error())
	}

	var ids []string
	for _, stream := range response.Data.Streams {
		ids = append(ids, stream.UserID)
	}

	return ids
}
