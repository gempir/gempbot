package helix

import (
	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
)

func (c *Client) GetTopChannels() []string {
	response, err := c.client.GetStreams(&helix.StreamsParams{Type: "live", First: 100})
	if err != nil {
		log.Error(err.Error())
	}

	var ids []string
	for _, stream := range response.Data.Streams {
		ids = append(ids, stream.UserID)
	}

	return ids
}