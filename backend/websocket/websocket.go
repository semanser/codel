package websocket

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	connections   map[string]*websocket.Conn // Map to store WebSocket connections
	connectionsMu sync.Mutex                 // Mutex to synchronize access to the connections map
)

func HandleWebsocket(c *gin.Context) {
	if connections == nil {
		connections = make(map[string]*websocket.Conn)
	}

	id := c.Param("id")

	// Upgrade HTTP connection to WebSocket
	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// Lock access to the connections map
	connectionsMu.Lock()
	defer connectionsMu.Unlock()

	// Save the connection in the map
	connections[id] = conn
}

func GetConnection(id string) (*websocket.Conn, error) {
	connectionsMu.Lock()
	defer connectionsMu.Unlock()
	conn, ok := connections[id]
	if !ok {
		return nil, fmt.Errorf("connection not found for id %s", id)
	}
	return conn, nil
}

func SendToChannel(id string, message string) error {
	conn, err := GetConnection(id)

	if err != nil {
		return fmt.Errorf("failed to send to the channel: %w", err)
	}

	return conn.WriteMessage(websocket.BinaryMessage, []byte(message))
}

func FormatTerminalInput(text string) string {
	yellow := "\033[33m" // ANSI escape code for yellow color
	reset := "\033[0m"   // ANSI escape code to reset color
	return fmt.Sprintf("$ %s%s%s\r\n", yellow, text, reset)
}

func FormatTerminalSystemOutput(text string) string {
	blue := "\033[34m" // ANSI escape code for blue color
	reset := "\033[0m" // ANSI escape code to reset color
	return fmt.Sprintf("%s%s%s\r\n", blue, text, reset)
}
