package ws

import (
	"encoding/json"
	"net/http"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/auth"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/media"
	"github.com/gorilla/websocket"
)

type WsHandler struct {
	upgrader     websocket.Upgrader
	authClient   *auth.Auth
	mediaManager *media.MediaManager
}

func NewWsHandler(authClient *auth.Auth, mediaManager *media.MediaManager) *WsHandler {
	return &WsHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		authClient:   authClient,
		mediaManager: mediaManager,
	}
}

type WsMessage struct {
	Message string `json:"message"`
}

func (h *WsHandler) HandleWs(w http.ResponseWriter, r *http.Request) {
	userId := ""
	if h.authClient.CanAuthenticate(r) {
		apiResp, _, apiErr := h.authClient.AttemptAuth(r, w)
		if apiErr != nil {
			api.WriteJson(w, "Auth error: "+apiErr.Error(), http.StatusUnauthorized)
			return
		}
		userId = apiResp.Data.UserID
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("ws upgrade failed: %s", err)
		return
	}

	writeQueue := make(chan []byte)

	connectionId := h.mediaManager.RegisterConnection(userId, func(message []byte) {
		writeQueue <- message
	})

	go startWriter(conn, writeQueue)

	defer func() {
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("ws read failed: %s", err)
			break
		}

		h.handleMessage(connectionId, userId, message)
	}
}

func startWriter(conn *websocket.Conn, writeQueue chan []byte) {
	for {
		message := <-writeQueue
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Errorf("Failed to write message %s", err.Error())
			return
		}
	}
}

type PlayerStateMessage struct {
	Action  string            `json:"action"`
	Seconds float32           `json:"seconds"`
	VideoId string            `json:"videoId"`
	State   media.PlayerState `json:"state"`
}

type Join struct {
	Action  string `json:"action"`
	Channel string `json:"channel"`
}

type BaseMessage struct {
	Action string `json:"action"`
}

func (h *WsHandler) handleMessage(connectionId string, userId string, byteMessage []byte) {
	var baseMessage BaseMessage
	err := json.Unmarshal(byteMessage, &baseMessage)
	if err != nil {
		log.Errorf("Failed to unmarshal message: %s", err)
		return
	}

	switch baseMessage.Action {
	case "PLAYER_STATE":
		var msg PlayerStateMessage
		err := json.Unmarshal(byteMessage, &msg)
		if err != nil {
			log.Errorf("Failed to unmarshal PlayerState message: %s", err)
			return
		}
		h.mediaManager.HandlePlayerState(connectionId, userId, msg.State, msg.VideoId, msg.Seconds)
	case "JOIN":
		var msg Join
		err := json.Unmarshal(byteMessage, &msg)
		if err != nil {
			log.Errorf("Failed to unmarshal Join message: %s", err)
			return
		}
		h.mediaManager.HandleJoin(connectionId, userId, msg.Channel)
	}
}
