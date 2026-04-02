package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now, rely on CORS/auth
	},
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients organized by Map ID.
	clients map[uuid.UUID]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan BroadcastMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	mu sync.RWMutex
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
	// The map this client is viewing
	mapID uuid.UUID
}

type BroadcastMessage struct {
	MapID   uuid.UUID
	Payload []byte
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uuid.UUID]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.mapID]; !ok {
				h.clients[client.mapID] = make(map[*Client]bool)
			}
			h.clients[client.mapID][client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.mapID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.mapID)
					}
				}
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			if clients, ok := h.clients[message.MapID]; ok {
				for client := range clients {
					select {
					case client.send <- message.Payload:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastMapUpdate(mapID uuid.UUID, updateType string, data interface{}) {
	payload := map[string]interface{}{
		"type": updateType,
		"data": data,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal broadcast message")
		return
	}
	h.broadcast <- BroadcastMessage{
		MapID:   mapID,
		Payload: bytes,
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("error reading websocket message")
			}
			break
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, authMiddleware authMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mapIDStr := chi.URLParam(r, "mapID")
		mapID, err := uuid.Parse(mapIDStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid map ID")
			return
		}

		// Auth via query parameter for WebSockets
		tokenStr := r.URL.Query().Get("token")
		if tokenStr == "" {
			respondError(w, http.StatusUnauthorized, "Missing token")
			return
		}
		
		userID, err := authMiddleware.validateToken(tokenStr)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		user, err := authMiddleware.userRepo.FindByID(userID)
		if err != nil || user == nil {
			respondError(w, http.StatusUnauthorized, "User not found")
			return
		}

		ctx := ctxWithUserID(r.Context(), user.ID.String())
		r = r.WithContext(ctx)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("Failed to upgrade to websocket")
			return
		}

		client := &Client{
			hub:   hub,
			conn:  conn,
			send:  make(chan []byte, 256),
			mapID: mapID,
		}

		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump()
	}
}