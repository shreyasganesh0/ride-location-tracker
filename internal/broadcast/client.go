package broadcast

import (
	"github.com/gorrila/websocket"
)

type Client struct {

	Hub *Hub
	Conn *websocket.Conn
	OutBoundMessagesCh chan []byte
}
