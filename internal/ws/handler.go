package ws

import (
	"encoding/json"
	"net/http"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/auth"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/media"
	"github.com/gorilla/websocket"
	"github.com/puzpuzpuz/xsync"
)

type WsHandler struct {
	upgrader     websocket.Upgrader
	authClient   *auth.Auth
	clients      *xsync.MapOf[string, *websocket.Conn]
	mediaManager *media.MediaManager
	writeQueues  *xsync.MapOf[string, chan []byte]
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
		clients:      xsync.NewMapOf[*websocket.Conn](),
		writeQueues:  xsync.NewMapOf[chan []byte](),
	}
}

type WsMessage struct {
	Message string `json:"message"`
}

func (h *WsHandler) HandleWs(w http.ResponseWriter, r *http.Request) {
	apiResp, _, apiErr := h.authClient.AttemptAuth(r, w)
	if apiErr != nil {
		api.WriteJson(w, "Auth error: "+apiErr.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("ws upgrade failed: %s", err)
		return
	}

	h.clients.Store(apiResp.Data.UserID, conn)
	writeQueue := make(chan []byte)
	h.writeQueues.Store(apiResp.Data.UserID, writeQueue)
	go startWriter(conn, writeQueue)
	h.mediaManager.RegisterWriter(apiResp.Data.UserID, func(message []byte) {
		writeQueue <- message
	})

	defer func() {
		h.clients.Delete(apiResp.Data.UserID)
		h.writeQueues.Delete(apiResp.Data.UserID)
		conn.Close()
	}()

	h.writeMessage(conn, WsMessage{"Authenticated you as " + apiResp.Data.Login})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("ws read failed: %s", err)
			break
		}
		h.handleMessage(apiResp.Data.UserID, message)
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

type TimeChanged struct {
	Action  string  `json:"action"`
	Seconds float32 `json:"seconds"`
	VideoId string  `json:"videoId"`
}

type BaseMessage struct {
	Action string `json:"action"`
}

func (h *WsHandler) handleMessage(userId string, byteMessage []byte) {
	var baseMessage BaseMessage
	err := json.Unmarshal(byteMessage, &baseMessage)
	if err != nil {
		log.Errorf("Failed to unmarshal message: %s", err)
		return
	}

	switch baseMessage.Action {
	case "TIME_CHANGED":
		var msg TimeChanged
		err := json.Unmarshal(byteMessage, &msg)
		if err != nil {
			log.Errorf("Failed to unmarshal TimeChanged message: %s", err)
			return
		}
		h.mediaManager.HandleTimeChange(userId, msg.VideoId, msg.Seconds)
	}
}

func (h *WsHandler) writeMessage(conn *websocket.Conn, message WsMessage) {
	if conn == nil {
		log.Error("Can't write on nil conn")
		return
	}

	result, err := json.Marshal(message)
	if err != nil {
		log.Error(err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, result)
	if err != nil {
		log.Error(err)
		return
	}
}
