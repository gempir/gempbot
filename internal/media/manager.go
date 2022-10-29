package media

import (
	"encoding/json"

	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/puzpuzpuz/xsync"
)

type MEDIA_TYPE string

const (
	MEDIA_TYPE_YOUTUBE MEDIA_TYPE = "youtube"
)

type MediaManager struct {
	db          store.Store
	rooms       *xsync.MapOf[string, *Room]
	connections *xsync.MapOf[string, *Connection]
}

type Connection struct {
	writer func(message []byte)
}

type Room struct {
	MediaType      MEDIA_TYPE
	CurrentVideoId string
	CurrentTime    float32
	users          *xsync.MapOf[string, *Connection]
}

func NewMediaManager(db store.Store) *MediaManager {
	return &MediaManager{
		db:          db,
		rooms:       xsync.NewMapOf[*Room](),
		connections: xsync.NewMapOf[*Connection](),
	}
}

func (m *MediaManager) HandleJoin(userID string, channel string) {
	joinChannelId := channel
	if channel == "" {
		joinChannelId = userID
	}

	connection, ok := m.connections.Load(userID)
	if !ok {
		return
	}

	room, ok := m.rooms.Load(joinChannelId)
	if !ok {
		room = &Room{users: xsync.NewMapOf[*Connection]()}
		m.rooms.Store(joinChannelId, room)
	}

	room.users.Store(userID, connection)
}

type TimeChangedMessage struct {
	Action      string  `json:"action"`
	VideoId     string  `json:"videoId"`
	CurrentTime float32 `json:"currentTime"`
}

func (m *MediaManager) HandleTimeChange(userID string, videoId string, currentTime float32) {
	state := m.getRoom(userID)

	state.CurrentTime = currentTime
	state.CurrentVideoId = videoId

	resultMessage, err := json.Marshal(TimeChangedMessage{CurrentTime: currentTime, VideoId: videoId, Action: "TIME_CHANGED"})
	if err != nil {
		log.Error(err)
	}

	state.users.Range(func(key string, conn *Connection) bool {
		conn.writer(resultMessage)
		return true
	})
}

func (m *MediaManager) getRoom(channelId string) *Room {
	state, ok := m.rooms.Load(channelId)
	if ok {
		return state
	}

	newState := &Room{MediaType: MEDIA_TYPE_YOUTUBE}
	m.rooms.Store(channelId, newState)
	return newState
}

func (m *MediaManager) RegisterConnection(userID string, writeFunc func(message []byte)) {
	m.connections.Store(userID, &Connection{writer: writeFunc})
}
