package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"live-polling-backend/services"
)

type ResultHandler struct {
	redisService *services.RedisService
}

func NewResultHandler(redisService *services.RedisService) *ResultHandler {
	return &ResultHandler{
		redisService: redisService,
	}
}

func (h *ResultHandler) GetResults(c *gin.Context) {

	// Get poll ID from URL
	pollID := c.Param("id")

	// Get current results from Redis
	results, err := h.redisService.GetVoteResults(
		c.Request.Context(),
		pollID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get poll results",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
	})
}
