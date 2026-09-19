package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"live-polling-backend/config"
	"live-polling-backend/handlers"
	"live-polling-backend/repositories"
	"live-polling-backend/routes"
	"live-polling-backend/services"
)

func main() {

	mongoClient, err := config.ConnectMongoDB()
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	redisClient, err := config.ConnectRedis()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	pollCollection := config.GetPollCollection(mongoClient)

	voteCollection := mongoClient.
		Database("live_polling").
		Collection("votes")

	userCollection := mongoClient.
		Database("live_polling").
		Collection("users")

	pollRepository := repositories.NewPollRepository(
		pollCollection,
	)

	pollService := services.NewPollService(
		pollRepository,
	)

	pollHandler := handlers.NewPollHandler(
		pollService,
	)

	redisService := services.NewRedisService(
		redisClient,
	)

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

	resultHandler := handlers.NewResultHandler(
		redisService,
	)

	webSocketService := services.NewWebSocketService()

	realtimeService := services.NewRealtimeService(
		redisClient,
		webSocketService,
	)

	webSocketHandler := handlers.NewWebSocketHandler(
		webSocketService,
		realtimeService,
	)

	userRepository := repositories.NewUserRepository(
		userCollection,
	)

	userService := services.NewUserService(
		userRepository,
	)

	authHandler := handlers.NewAuthHandler(
		userService,
	)

	router := gin.Default()

	frontendURL := os.Getenv("FRONTEND_URL")

	allowOrigins := []string{
		"http://localhost:5173",
	}

	if frontendURL != "" {
		allowOrigins = append(
			allowOrigins,
			frontendURL,
		)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Live Polling Backend is running!",
		})
	})

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

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	err = router.Run(":" + port)

	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
