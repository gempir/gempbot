package media

import (
	"github.com/gempir/gempbot/internal/store"
)

type MEDIA_TYPE string

const (
	MEDIA_TYPE_YOUTUBE MEDIA_TYPE = "youtube"
)

type MediaManager struct {
	db     store.Store
	states map[string]*MediaPlayerState
}

type MediaPlayerState struct {
	MediaType      MEDIA_TYPE
	CurrentVideoId string
	CurrentTime    float32
}

func NewMediaManager(db store.Store) *MediaManager {
	return &MediaManager{
		db:     db,
		states: make(map[string]*MediaPlayerState),
	}
}

func (m *MediaManager) HandleTimeChange(channelId string, videoId string, currentTime float32) {
	state := m.getState(channelId)

	state.CurrentTime = currentTime
	state.CurrentVideoId = videoId
}

func (m *MediaManager) getState(channelId string) *MediaPlayerState {
	state, ok := m.states[channelId]
	if ok {
		return state
	}

	newState := &MediaPlayerState{MediaType: MEDIA_TYPE_YOUTUBE}
	m.states[channelId] = newState
	return newState
}
