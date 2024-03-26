package websocket

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	connections map[int64]*websocket.Conn // Map to store WebSocket connections
)

func HandleWebsocket(c *gin.Context) {
	if connections == nil {
		connections = make(map[int64]*websocket.Conn)
	}

	id := c.Param("id")

	parsedID, err := strconv.ParseInt(id, 10, 32)

	if err != nil {
		c.AbortWithError(400, fmt.Errorf("failed to parse id: %w", err))
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// Save the connection in the map
	connections[parsedID] = conn
}

func GetConnection(id int64) (*websocket.Conn, error) {
	conn, ok := connections[id]
	if !ok {
		return nil, fmt.Errorf("connection not found for id %d", id)
	}
	return conn, nil
}

func SendToChannel(id int64, message string) error {
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
