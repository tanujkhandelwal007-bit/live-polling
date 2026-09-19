package routes

import (
	"github.com/gin-gonic/gin"

	"live-polling-backend/handlers"
	"live-polling-backend/middleware"
)

func RegisterPollRoutes(
	router *gin.Engine,
	pollHandler *handlers.PollHandler,
) {
	// Creating a poll requires authentication.
	router.POST(
		"/polls",
		middleware.AuthMiddleware(),
		pollHandler.CreatePoll,
	)

	// Viewing a poll remains public.
	router.GET("/polls/:id", pollHandler.GetPollByID)

	// Closing a poll requires authentication.
	router.PUT(
		"/polls/:id/close",
		middleware.AuthMiddleware(),
		pollHandler.ClosePoll,
	)
}
