package sync

import "github.com/gorilla/websocket"

// TODO: Phase 4 - Implement WebSocket sync client

// Client manages sync with remote peers
type Client struct {
   ws *websocket.Conn
}

// NewClient creates a new sync client
func NewClient(relayURL string) *Client {
	return &Client{	ws: nil}
}

// Connect establishes connection to relay server
func (c *Client) Connect() error {
	return nil
}

// SendOperation sends an operation to peers
func (c *Client) SendOperation(opData []byte) error {
	c.ws.WriteMessage(websocket.BinaryMessage, opData)
	return nil
}
