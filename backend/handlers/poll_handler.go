package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"live-polling-backend/services"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PollHandler struct {
	service *services.PollService
}

func NewPollHandler(service *services.PollService) *PollHandler {
	return &PollHandler{
		service: service,
	}
}

type CreatePollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

func (h *PollHandler) CreatePoll(c *gin.Context) {

	var request CreatePollRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// Get logged-in user's ID from authentication middleware.
	userID, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user authentication required",
		})
		return
	}

	createdBy, ok := userID.(string)

	if !ok || createdBy == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	poll, err := h.service.CreatePoll(
		c.Request.Context(),
		request.Question,
		request.Options,
		createdBy,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "poll created successfully",
		"poll":    poll,
	})
}

func (h *PollHandler) GetPollByID(c *gin.Context) {

	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid poll ID",
		})
		return
	}

	poll, err := h.service.GetPollByID(
		c.Request.Context(),
		objectID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "poll not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"poll": poll,
	})
}

func (h *PollHandler) ClosePoll(c *gin.Context) {

	id := c.Param("id")

	objectID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid poll ID",
		})
		return
	}

	// Get logged-in user's ID from authentication middleware.
	userID, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user authentication required",
		})
		return
	}

	userIDString, ok := userID.(string)

	if !ok || userIDString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	err = h.service.ClosePoll(
		c.Request.Context(),
		objectID,
		userIDString,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "poll closed successfully",
	})
}
