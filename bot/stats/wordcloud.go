package stats

import (
	"strings"
)

var (
	ignoredWords = []string{"a", "the", "you", "i", "is", "to", "no", "yes"}
)

func (b *Broadcaster) worldcloudMessageHandler() {
	go func() {
		for message := range b.messageQueue {

			words := []string{}

			for _, word := range strings.Split(message.Message, " ") {
				if !isIgnored(word) {
					words = append(words, word)
				}
			}

			b.store.AddToWordcloud(words...)
		}
	}()
}

func isIgnored(word string) bool {
	for _, ignored := range ignoredWords {
		if strings.ToLower(word) == strings.ToLower(ignored) {
			return true
		}
	}
	return false
}
