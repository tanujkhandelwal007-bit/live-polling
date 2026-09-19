package routes

import (
	"github.com/gin-gonic/gin"

	"live-polling-backend/handlers"
)

func RegisterWebSocketRoutes(
	router *gin.Engine,
	webSocketHandler *handlers.WebSocketHandler,
) {
	router.GET(
		"/polls/:id/ws",
		webSocketHandler.HandleConnection,
	)
}
