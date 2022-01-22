package channelpoint

import (
	"encoding/json"

	"github.com/gempir/gempbot/internal/log"
)

func UnmarshallBttvAdditionalOptions(jsonString string) BttvAdditionalOptions {
	if jsonString == "{}" {
		return BttvAdditionalOptions{Slots: 1}
	}

	var additionalOptions BttvAdditionalOptions

	if err := json.Unmarshal([]byte(jsonString), &additionalOptions); err != nil {
		log.Error(err)
		return BttvAdditionalOptions{Slots: 1}
	}

	return additionalOptions
}

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
