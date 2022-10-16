package ws

import (
	"encoding/json"
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
	defer conn.Close()

	h.writeMessage(conn, WsMessage{"Authenticated you as " + apiResp.Data.Login})

	// Continuosly read and write message
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("ws read failed: %s", err)
			break
		}
		input := string(message)
		log.Infof("recv: %d %s", mt, input)
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
