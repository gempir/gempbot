package media

import (
	"github.com/gempir/gempbot/internal/store"
	"github.com/puzpuzpuz/xsync"
)

type MEDIA_TYPE string

const (
	MEDIA_TYPE_YOUTUBE MEDIA_TYPE = "youtube"
)

type MediaManager struct {
	db      store.Store
	states  *xsync.MapOf[string, *MediaPlayerState]
	writers *xsync.MapOf[string, func(message []byte)]
}

type MediaPlayerState struct {
	MediaType      MEDIA_TYPE
	CurrentVideoId string
	CurrentTime    float32
}

func NewMediaManager(db store.Store) *MediaManager {
	return &MediaManager{
		db:      db,
		states:  xsync.NewMapOf[*MediaPlayerState](),
		writers: xsync.NewMapOf[func(message []byte)](),
	}
}

func (m *MediaManager) HandleTimeChange(userID string, videoId string, currentTime float32) {
	state := m.getState(userID)

	state.CurrentTime = currentTime
	state.CurrentVideoId = videoId

	// writer, ok := m.writers.Load(userID)
	// if !ok {
	// 	return
	// }
	// writer([]byte("time change"))
}

func (m *MediaManager) getState(channelId string) *MediaPlayerState {
	state, ok := m.states.Load(channelId)
	if ok {
		return state
	}

	newState := &MediaPlayerState{MediaType: MEDIA_TYPE_YOUTUBE}
	m.states.Store(channelId, newState)
	return newState
}

func (m *MediaManager) RegisterWriter(userID string, writeFunc func(message []byte)) {
	m.writers.Store(userID, writeFunc)
}
