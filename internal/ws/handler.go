package ws

import (
	"fmt"
	"net/http"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/auth"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gorilla/websocket"
)

type WsHandler struct {
	upgrader   websocket.Upgrader
	authClient *auth.Auth
}

func NewWsHandler(authClient *auth.Auth) *WsHandler {
	return &WsHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		authClient: authClient,
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

	go startWriter(conn, writeQueue)

	defer func() {
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("Client Disconnected %s %s", userId, err)
			break
		}

		fmt.Println(message)
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
