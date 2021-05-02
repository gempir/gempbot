package server

import (
	"net/http"
	"sort"
	"sync"

	"github.com/gempir/spamchamp/services/api/emotechief"

	"github.com/gempir/spamchamp/pkg/config"
	"github.com/gempir/spamchamp/pkg/helix"
	"github.com/gempir/spamchamp/pkg/store"
	"github.com/rs/cors"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// Server api server
type Server struct {
	cfg             *config.Config
	broadcastQueue  chan BroadcastMessage
	helixClient     *helix.Client
	helixUserClient *helix.Client
	store           *store.Store
	emotechief      *emotechief.EmoteChief
}

type BroadcastMessage struct {
	Records        []Record `json:"records"`
	JoinedChannels int      `json:"joinedChannels"`
	ActiveChannels int      `json:"activeChannels"`
}

type WordcloudWord struct {
	Text  string  `json:"text"`
	Value float64 `json:"value"`
}

type Record struct {
	Title  string  `json:"title"`
	Scores []Score `json:"scores"`
}

type Score struct {
	User  User    `json:"user"`
	Score float64 `json:"score"`
}

type User struct {
	Id             string `json:"id"`
	DisplayName    string `json:"displayName"`
	ProfilePicture string `json:"profilePicture"`
}

func (s *Record) GetScoresSorted() []Score {
	scores := s.Scores
	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})
	return scores
}

// NewServer create api Server
func NewServer(cfg *config.Config, helixClient *helix.Client, helixUserClient *helix.Client, store *store.Store, emotechief *emotechief.EmoteChief, broadcastQueue chan BroadcastMessage) Server {
	return Server{
		cfg:             cfg,
		broadcastQueue:  broadcastQueue,
		helixClient:     helixClient,
		helixUserClient: helixUserClient,
		store:           store,
		emotechief:      emotechief,
	}
}

// var clients = make(map[*websocket.Conn]bool) // connected clients
var clients sync.Map
var upgrader = websocket.Upgrader{}

func (s *Server) Start() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	go s.handleMessages()
	go s.syncSubscriptions()
	go s.tokenRefreshRoutine()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/ws", s.handleConnections)
	mux.HandleFunc("/api/callback", s.handleCallback)
	mux.HandleFunc("/api/redemption", s.handleChannelPointsRedemption)
	mux.HandleFunc("/api/userConfig", s.handleUserConfig)
	mux.HandleFunc("/api/rewards", s.handleRewards)

	handler := cors.AllowAll().Handler(mux)
	log.Info("listening on port :8035")
	err := http.ListenAndServe(":8035", handler)
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	clients.Store(ws, true)
	// Make sure we close the connection when the function returns
	defer ws.Close()

	for {
		var msg BroadcastMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Debugf("handleConnection error: %v", err)
			clients.Delete(ws)
			break
		}
		// Send the newly received message to the broadcast channel
		s.broadcastQueue <- msg
	}
}

func (s *Server) handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-s.broadcastQueue
		// Send it out to every client that is currently connected
		clients.Range(func(client, value interface{}) bool {
			res := client.(*websocket.Conn)
			err := res.WriteJSON(msg)
			if err != nil {
				log.Errorf("broadcast error: %v", err)
				res.Close()
				clients.Delete(client)
			}

			return true
		})

	}
}
