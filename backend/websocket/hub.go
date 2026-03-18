package websocket

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uuid.UUID]map[*Client]bool),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[c.userID]; !ok {
		h.clients[c.userID] = make(map[*Client]bool)
	}
	h.clients[c.userID][c] = true
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	userClients, ok := h.clients[c.userID]
	if !ok {
		return
	}

	if _, ok := userClients[c]; ok {
		delete(userClients, c)
		close(c.send)
	}

	if len(userClients) == 0 {
		delete(h.clients, c.userID)
	}
}

func (h *Hub) BroadcastToUser(userID uuid.UUID, msg Message) {
	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mu.RLock()
	userClients := h.clients[userID]
	h.mu.RUnlock()

	for c := range userClients {
		select {
		case c.send <- payload:
		default:
			// Drop if client is slow; it will reconnect.
		}
	}
}
