package routes

import (
	"github.com/gin-gonic/gin"

	"live-polling-backend/handlers"
)

func RegisterAuthRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
) {
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
}
