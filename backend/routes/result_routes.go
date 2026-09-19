package routes

import (
	"github.com/gin-gonic/gin"

	"live-polling-backend/handlers"
)

func RegisterResultRoutes(
	router *gin.Engine,
	resultHandler *handlers.ResultHandler,
) {
	router.GET("/polls/:id/results", resultHandler.GetResults)
}
