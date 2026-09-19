# Live Polling Backend

A Go backend for creating polls, collecting votes, showing live results, and broadcasting realtime updates.

## Project Structure

```text
backend/
├── main.go                    # Application entry point and dependency wiring
├── go.mod                     # Go module and dependency definitions
├── README.md                  # Project documentation
│
├── config/                    # External service connections
│   ├── mongo.go               # MongoDB connection and poll collection setup
│   └── redis.go               # Redis connection setup
│
├── handlers/                 # HTTP and WebSocket request handlers
│   ├── auth_handler.go        # Registration and login handlers
│   ├── poll_handler.go        # Poll creation, reading, and closing handlers
│   ├── result_handler.go      # Poll result handler
│   ├── vote_handler.go        # Vote submission handler
│   └── websocket_handler.go   # WebSocket connection handler
│
├── middleware/               # Request middleware
│   └── auth.go                # JWT Bearer-token authentication
│
├── models/                   # MongoDB data models
│   ├── poll.go                # Poll model
│   ├── user.go                # User model
│   └── vote.go                # Vote model
│
├── repositories/             # Database access layer
│   ├── poll_repository.go     # Poll database operations
│   ├── user_repository.go     # User database operations
│   └── vote_repository.go     # Vote database operations
│
├── routes/                   # API route registration
│   ├── auth_routes.go         # Authentication routes
│   ├── poll_routes.go         # Poll routes
│   ├── result_routes.go       # Result routes
│   ├── vote_routes.go         # Voting routes
│   └── websocket_routes.go    # WebSocket routes
│
└── services/                 # Application business logic
    ├── poll_service.go        # Poll validation and lifecycle logic
    ├── realtime_service.go    # Redis Pub/Sub listener
    ├── redis_service.go       # Live vote counts and event publishing
    ├── user_service.go        # Registration, login, hashing, and JWT creation
    ├── vote_service.go        # Vote validation and persistence
    └── websocket_service.go   # Connected client management and broadcasting
```

## Main Technologies

- Go
- Gin
- MongoDB
- Redis
- JWT authentication
- bcrypt password hashing
- WebSockets

## Services Used

The backend expects these services to be running locally:

```text
MongoDB: localhost:27017
Redis:   localhost:6379
API:     localhost:8080
```
