package channelpoint

import (
	"encoding/json"

	"github.com/gempir/gempbot/internal/log"
)

func UnmarshallSevenTvAdditionalOptions(jsonString string) SevenTvAdditionalOptions {
	if jsonString == "{}" {
		return SevenTvAdditionalOptions{Slots: 1}
	}

	var additionalOptions SevenTvAdditionalOptions

	if err := json.Unmarshal([]byte(jsonString), &additionalOptions); err != nil {
		log.Error(err)
		return SevenTvAdditionalOptions{Slots: 1}
	}

	return additionalOptions
}
