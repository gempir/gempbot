package media

import (
	"encoding/json"

	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync"
)

type PlayerState int

const (
	UNSTARTED PlayerState = -1
	ENDED     PlayerState = 0
	PLAYING   PlayerState = 1
	PAUSED    PlayerState = 2
	BUFFERING PlayerState = 3
	CUED      PlayerState = 5
)

type MEDIA_TYPE string

const (
	MEDIA_TYPE_YOUTUBE MEDIA_TYPE = "youtube"
)

type DebugMessage struct {
	Action  string `json:"string"`
	Message string `json:"message"`
}

type MediaManager struct {
	db          store.Store
	helixClient helixclient.Client
	rooms       *xsync.MapOf[string, *Room]
	connections *xsync.MapOf[string, *Connection]
}

type Connection struct {
	id     string
	writer func(message []byte)
}

type Room struct {
	MediaType      MEDIA_TYPE
	CurrentVideoId string
	CurrentTime    float32
	State          PlayerState
	users          *xsync.MapOf[string, *Connection]
}

func NewMediaManager(db store.Store, helixClient helixclient.Client) *MediaManager {
	return &MediaManager{
		db:          db,
		helixClient: helixClient,
		rooms:       xsync.NewMapOf[*Room](),
		connections: xsync.NewMapOf[*Connection](),
	}
}

func (m *MediaManager) HandleJoin(connectionId string, userID string, channel string) {
	var joinChannelId string
	if channel == "" {
		joinChannelId = userID
	} else {
		res, err := m.helixClient.GetUserByUsername(channel)
		if err != nil {
			return
		}
		joinChannelId = res.ID
	}

	connection, ok := m.connections.Load(connectionId)
	if !ok {
		return
	}

	room, ok := m.rooms.Load(joinChannelId)
	if !ok {
		room = newRoom()
		m.rooms.Store(joinChannelId, room)
	}

	sendPlayerState([]*Connection{connection}, room)

	room.users.Store(connectionId, connection)
}

type PlayerStateMessage struct {
	Action      string      `json:"action"`
	VideoId     string      `json:"videoId"`
	CurrentTime float32     `json:"currentTime"`
	State       PlayerState `json:"state"`
}

func (m *MediaManager) HandlePlayerState(connectionId string, userID string, state PlayerState, videoId string, currentTime float32) {
	if userID == "" {
		log.Errorf("missing userID time %f on connection %s", currentTime, connectionId)
		return
	}

	roomState := m.getRoom(userID)

	roomState.CurrentTime = currentTime
	roomState.CurrentVideoId = videoId
	roomState.State = state

	conns := []*Connection{}
	roomState.users.Range(func(key string, conn *Connection) bool {
		if conn.id != connectionId {
			conns = append(conns, conn)
		}
		return true
	})

	if roomState.CurrentVideoId != "" {
		sendPlayerState(conns, roomState)
	}
}

func (m *MediaManager) getRoom(channelId string) *Room {
	room, ok := m.rooms.Load(channelId)
	if ok {
		return room
	}

	newRoom := newRoom()
	m.rooms.Store(channelId, newRoom)

	return newRoom
}

func (m *MediaManager) RegisterConnection(userID string, writeFunc func(message []byte)) string {
	connectionId := uuid.NewString()

	m.connections.Store(connectionId, &Connection{writer: writeFunc, id: connectionId})

	return connectionId
}

func newRoom() *Room {
	return &Room{
		users:     xsync.NewMapOf[*Connection](),
		MediaType: MEDIA_TYPE_YOUTUBE,
	}
}

func sendPlayerState(connections []*Connection, room *Room) {
	resultMessage, err := json.Marshal(newPlayerStateMessage(room))
	if err != nil {
		log.Error(err)
		return
	}

	for _, conn := range connections {
		conn.writer(resultMessage)
	}
}

func newPlayerStateMessage(room *Room) PlayerStateMessage {
	return PlayerStateMessage{
		Action:      "PLAYER_STATE",
		VideoId:     room.CurrentVideoId,
		CurrentTime: room.CurrentTime,
		State:       room.State,
	}
}
