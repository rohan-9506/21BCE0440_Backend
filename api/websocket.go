package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade the HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocketHub manages WebSocket connections and broadcasts messages
type WebSocketHub struct {
	clients map[string]*websocket.Conn
}

// NewWebSocketHub creates a new WebSocketHub
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients: make(map[string]*websocket.Conn),
	}
}

// AddClient adds a new client to the hub
func (hub *WebSocketHub) AddClient(userID string, conn *websocket.Conn) {
	hub.clients[userID] = conn
}

// RemoveClient removes a client from the hub
func (hub *WebSocketHub) RemoveClient(userID string) {
	delete(hub.clients, userID)
}

// NotifyFileUpload sends a notification message to the WebSocket connection
func (hub *WebSocketHub) NotifyFileUpload(userID, fileURL string) {
	if conn, ok := hub.clients[userID]; ok {
		err := conn.WriteMessage(websocket.TextMessage, []byte(fileURL))
		if err != nil {
			hub.RemoveClient(userID)
		}
	}
}

var hub = NewWebSocketHub()

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(c *gin.Context) {
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set up WebSocket connection"})
		return
	}
	defer conn.Close()

	hub.AddClient(userID, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			hub.RemoveClient(userID)
			break
		}
	}
}
