package services

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

type WebSocketService struct {
	clients map[string]map[*WebSocketClient]bool
	mutex   sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		clients: make(map[string]map[*WebSocketClient]bool),
	}
}

// AddClient adds a browser connection to a poll.
func (s *WebSocketService) AddClient(
	pollID string,
	conn *websocket.Conn,
) *WebSocketClient {

	client := &WebSocketClient{
		conn: conn,
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.clients[pollID] == nil {
		s.clients[pollID] = make(map[*WebSocketClient]bool)
	}

	s.clients[pollID][client] = true

	return client
}

// RemoveClient removes a browser connection from a poll.
func (s *WebSocketService) RemoveClient(
	pollID string,
	client *WebSocketClient,
) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.clients[pollID] != nil {
		delete(s.clients[pollID], client)
	}

	_ = client.conn.Close()
}

// Broadcast sends a message to all clients connected to a poll.
func (s *WebSocketService) Broadcast(
	pollID string,
	message []byte,
) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for client := range s.clients[pollID] {

		// Lock this particular WebSocket connection.
		client.mutex.Lock()

		err := client.conn.WriteMessage(
			websocket.TextMessage,
			message,
		)

		client.mutex.Unlock()

		if err != nil {
			// Ignore failed connections for now.
			continue
		}
	}
}
