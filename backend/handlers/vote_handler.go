package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"live-polling-backend/services"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type VoteHandler struct {
	service *services.VoteService
}

func NewVoteHandler(service *services.VoteService) *VoteHandler {
	return &VoteHandler{
		service: service,
	}
}

type CreateVoteRequest struct {
	Option string `json:"option"`
}

func (h *VoteHandler) CreateVote(c *gin.Context) {

	// Get poll ID from URL
	pollID := c.Param("id")

	// Convert string ID to MongoDB ObjectID
	objectID, err := bson.ObjectIDFromHex(pollID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid poll ID",
		})
		return
	}

	// Read JSON body
	var request CreateVoteRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// Create vote
	vote, err := h.service.CreateVote(
		c.Request.Context(),
		objectID,
		request.Option,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "vote created successfully",
		"vote":    vote,
	})
}
