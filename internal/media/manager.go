package media

import (
	"encoding/json"
	"regexp"

	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync"
)

type PlayerState string

const (
	PLAYING PlayerState = "PLAYING"
	PAUSED  PlayerState = "PAUSED"
)

var (
	YOUTUBE_REGEX = regexp.MustCompile(`^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)
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
	storage                   storage
	helixClient               helixclient.Client
	rooms                     *xsync.MapOf[string, *Room]
	connections               *xsync.MapOf[string, *Connection]
	bot                       mediaBot
	commandsActivatedChannels map[string]bool
}

type Connection struct {
	id     string
	writer func(message []byte)
}

type Room struct {
	MediaType MEDIA_TYPE
	Url       string
	Time      float32
	State     PlayerState
	users     *xsync.MapOf[string, *Connection]
}

type storage interface {
	AddToQueue(queueItem store.MediaQueue) error
	GetAllMediaCommandsBotConfig() []store.BotConfig
}

type mediaBot interface {
	RegisterCommand(command string, handler func(dto.CommandPayload))
	Say(channel string, message string)
	Reply(channel string, parentMsgId, message string)
}

func NewMediaManager(storage storage, helixClient helixclient.Client, bot mediaBot) *MediaManager {

	commandsActivatedChannels := make(map[string]bool)
	commandActivatedCfgs := storage.GetAllMediaCommandsBotConfig()
	for _, cfg := range commandActivatedCfgs {
		if cfg.MediaCommands {
			commandsActivatedChannels[cfg.OwnerTwitchID] = true
		}
	}

	mm := &MediaManager{
		storage:                   storage,
		helixClient:               helixClient,
		rooms:                     xsync.NewMapOf[*Room](),
		connections:               xsync.NewMapOf[*Connection](),
		commandsActivatedChannels: commandsActivatedChannels,
		bot:                       bot,
	}

	bot.RegisterCommand("sr", mm.handleSongRequest)

	return mm
}

func (m *MediaManager) handleSongRequest(payload dto.CommandPayload) {
	if _, ok := m.commandsActivatedChannels[payload.Msg.RoomID]; !ok {
		return
	}

	if !YOUTUBE_REGEX.MatchString(payload.Query) {
		m.bot.Reply(payload.Msg.Channel, payload.Msg.ID, "invalid youtube url")
		return
	}

	m.AddUrlToQueue(payload.Query, payload.Msg.User.ID, payload.Msg.RoomID)
}

func (m *MediaManager) AddUrlToQueue(url string, authorID string, channelID string) {
	m.storage.AddToQueue(store.MediaQueue{
		ChannelTwitchId: channelID,
		Author:          authorID,
		Url:             url,
	})
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
	Action string      `json:"action"`
	Url    string      `json:"url"`
	Time   float32     `json:"time"`
	State  PlayerState `json:"state"`
}

func (m *MediaManager) HandlePlayerState(connectionId string, userID string, state PlayerState, url string, time float32) {
	if userID == "" {
		log.Errorf("missing userID time %f on connection %s", time, connectionId)
		return
	}

	roomState := m.getRoom(userID)

	roomState.Time = time
	roomState.Url = url
	roomState.State = state

	conns := []*Connection{}
	roomState.users.Range(func(key string, conn *Connection) bool {
		if conn.id != connectionId {
			conns = append(conns, conn)
		}
		return true
	})

	if roomState.Url != "" {
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

func (m *MediaManager) DeregisterConnection(connectionId string) {
	m.connections.Delete(connectionId)
}

func newRoom() *Room {
	return &Room{
		users: xsync.NewMapOf[*Connection](),
		Url:   "https://www.youtube.com/watch?v=wzE2nsjsHhg",
		Time:  0,
		State: PAUSED,
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
		Action: "PLAYER_STATE",
		Url:    room.Url,
		Time:   room.Time,
		State:  room.State,
	}
}
