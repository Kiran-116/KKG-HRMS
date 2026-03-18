package websocket

import (
	"net/http"

	"hrms/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// CORS is handled elsewhere; allow dev usage.
		return true
	},
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	// Auth via query param token (browser-friendly)
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "token is required"})
		return
	}

	claims, err := utils.ValidateToken(token, false)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Invalid or expired token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(h.hub, conn, claims.UserID)
	h.hub.Register(client)

	go client.writePump()
	go client.readPump()
}

