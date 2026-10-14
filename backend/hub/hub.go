package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	connections map[int64][]*websocket.Conn
	mu          sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[int64][]*websocket.Conn),
	}
}

func (h *Hub) Register(userID int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.connections[userID] = append(h.connections[userID], conn)
}

func (h *Hub) Unregister(userID int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns := h.connections[userID]

	for i, c := range conns {
		if c == conn {
			h.connections[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
}

func (h *Hub) Broadcast(userID int64, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conns := h.connections[userID]
	for _, c := range conns {
		c.WriteMessage(websocket.TextMessage, message)
	}
}
