package websocket

import (
	"encoding/json"
	"sync"

	"github.com/Micah-Shallom/geoint-backend/internal/models"
	"github.com/Micah-Shallom/geoint-backend/utility"
	"github.com/gorilla/websocket"
)

type Client struct {
	Hub        *Hub
	Conn       *websocket.Conn
	Send       chan []byte
	AnalysisID string
}

type Hub struct {
	Clients    map[string]map[*Client]bool // analysisID -> clients
	Broadcast  chan *models.Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
	Logger     *utility.Logger
}

func NewHub(logger *utility.Logger) *Hub {
	return &Hub{
		Clients:    make(map[string]map[*Client]bool),
		Broadcast:  make(chan *models.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Clients[client.AnalysisID] == nil {
				h.Clients[client.AnalysisID] = make(map[*Client]bool)
			}
			h.Clients[client.AnalysisID][client] = true
			h.mu.Unlock()
			h.Logger.Info("websocket client has been registered", client.AnalysisID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Clients[client.AnalysisID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.Clients, client.AnalysisID)
					}
				}
			}
			h.mu.Unlock()
			h.Logger.Info("WebSocket client unregistered", client.AnalysisID)

		case message := <-h.Broadcast:
			h.mu.RLock()
			clients := h.Clients[message.AnalysisID]
			h.mu.RUnlock()

			messageBytes, err := json.Marshal(message)
			if err != nil {
				h.Logger.Error("failed to marshal message", err)
				continue
			}

			for client := range clients {
				select {
				case client.Send <- messageBytes:
				default:
					h.mu.Lock()
					close(client.Send)
					delete(h.Clients[message.AnalysisID], client)
					h.mu.Unlock()
				}
			}
		}
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
