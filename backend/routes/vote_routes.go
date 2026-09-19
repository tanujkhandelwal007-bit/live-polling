package routes

import (
	"live-polling-backend/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterVoteRoutes(
	router *gin.Engine,
	voteHandler *handlers.VoteHandler,
) {
	router.POST("/polls/:id/vote", voteHandler.CreateVote)
}
