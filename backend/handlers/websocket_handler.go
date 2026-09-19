package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"live-polling-backend/services"
)

type WebSocketHandler struct {
	webSocketService *services.WebSocketService
	realtimeService  *services.RealtimeService
}

func NewWebSocketHandler(
	webSocketService *services.WebSocketService,
	realtimeService *services.RealtimeService,
) *WebSocketHandler {
	return &WebSocketHandler{
		webSocketService: webSocketService,
		realtimeService:  realtimeService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {

	// Get poll ID from URL
	pollID := c.Param("id")

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)

	if err != nil {
		return
	}

	// Add client to this poll
	client := h.webSocketService.AddClient(
		pollID,
		conn,
	)

	// Start Redis listener for this poll
	go h.realtimeService.ListenToPoll(
		c.Request.Context(),
		pollID,
	)

	// Remove client when connection closes
	defer h.webSocketService.RemoveClient(
		pollID,
		client,
	)

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()

		if err != nil {
			break
		}
	}
}
