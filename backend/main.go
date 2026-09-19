package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"live-polling-backend/config"
	"live-polling-backend/handlers"
	"live-polling-backend/repositories"
	"live-polling-backend/routes"
	"live-polling-backend/services"
)

func main() {

	// Connect to MongoDB
	mongoClient, err := config.ConnectMongoDB()
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	// Connect to Redis
	redisClient, err := config.ConnectRedis()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	// Get MongoDB collections
	pollCollection := config.GetPollCollection(mongoClient)

	voteCollection := mongoClient.
		Database("live_polling").
		Collection("votes")

	userCollection := mongoClient.
		Database("live_polling").
		Collection("users")

	// -------------------------
	// Poll dependencies
	// -------------------------

	pollRepository := repositories.NewPollRepository(
		pollCollection,
	)

	pollService := services.NewPollService(
		pollRepository,
	)

	pollHandler := handlers.NewPollHandler(
		pollService,
	)

	// -------------------------
	// Redis service
	// -------------------------

	redisService := services.NewRedisService(
		redisClient,
	)

	// -------------------------
	// Vote dependencies
	// -------------------------

	voteRepository := repositories.NewVoteRepository(
		voteCollection,
	)

	voteService := services.NewVoteService(
		voteRepository,
		pollRepository,
		redisService,
	)

	voteHandler := handlers.NewVoteHandler(
		voteService,
	)

	// -------------------------
	// Result handler
	// -------------------------

	resultHandler := handlers.NewResultHandler(
		redisService,
	)

	// -------------------------
	// WebSocket service
	// -------------------------

	webSocketService := services.NewWebSocketService()

	// -------------------------
	// Realtime service
	// -------------------------

	realtimeService := services.NewRealtimeService(
		redisClient,
		webSocketService,
	)

	// -------------------------
	// WebSocket handler
	// -------------------------

	webSocketHandler := handlers.NewWebSocketHandler(
		webSocketService,
		realtimeService,
	)

	// -------------------------
	// User dependencies
	// -------------------------

	userRepository := repositories.NewUserRepository(
		userCollection,
	)

	userService := services.NewUserService(
		userRepository,
	)

	authHandler := handlers.NewAuthHandler(
		userService,
	)

	// -------------------------
	// Gin router
	// -------------------------

	router := gin.Default()

	// -------------------------
	// CORS
	// -------------------------

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// -------------------------
	// Test endpoint
	// -------------------------

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Live Polling Backend is running!",
		})
	})

	// -------------------------
	// Register routes
	// -------------------------

	routes.RegisterPollRoutes(
		router,
		pollHandler,
	)

	routes.RegisterVoteRoutes(
		router,
		voteHandler,
	)

	routes.RegisterResultRoutes(
		router,
		resultHandler,
	)

	routes.RegisterWebSocketRoutes(
		router,
		webSocketHandler,
	)

	routes.RegisterAuthRoutes(
		router,
		authHandler,
	)

	// -------------------------
	// Start server
	// -------------------------

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
