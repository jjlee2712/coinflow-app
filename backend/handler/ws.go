package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jjlee2712/coinflow/backend/hub"
	"github.com/jjlee2712/coinflow/backend/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	Hub       *hub.Hub
	JWTSecret string
}

func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	userID, err := service.ValidateToken(token, h.JWTSecret)
	if err != nil {
		http.Error(w, "Error", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.Hub.Register(userID, conn)
	defer h.Hub.Unregister(userID, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
