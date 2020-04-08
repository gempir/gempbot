package stats

import (
	"strings"
)

func (b *Broadcaster) worldcloudMessageHandler() {
	go func() {
		for message := range b.messageQueue {
			b.store.AddToWordcloud(strings.Split(message.Message, " ")...)
		}
	}()
}
