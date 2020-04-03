package api

import (
	"github.com/gempir/spamchamp/bot/config"
	"github.com/gempir/spamchamp/bot/helix"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// Server api server
type Server struct {
	cfg            *config.Config
	broadcastQueue chan BroadcastMessage
	helixClient    *helix.Client
}

type BroadcastMessage struct {
	Channels map[string]FrontendStats `json:"channels"`
}

type FrontendStats struct {
	ChannelName       string `json:"channelName"`
	MessagesPerSecond int    `json:"messagesPerSecond"`
}

// NewServer create api Server
func NewServer(cfg *config.Config, helixClient *helix.Client, broadcastQueue chan BroadcastMessage) Server {
	return Server{
		cfg:            cfg,
		broadcastQueue: broadcastQueue,
		helixClient:    helixClient,
	}
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var upgrader = websocket.Upgrader{}

func (s *Server) Start() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	go s.handleMessages()
	http.HandleFunc("/api/ws", s.handleConnections)
	http.HandleFunc("/api/channels", s.handleChannels)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (s *Server) handleChannels(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	clients[ws] = true
	// Make sure we close the connection when the function returns
	defer ws.Close()

	for {
		var msg BroadcastMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Infof("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		s.broadcastQueue <- msg
	}
}

func (s *Server) handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <- s.broadcastQueue
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Errorf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
