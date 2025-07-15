package broadcast

import (
	"github.com/gorrila/websocket"
)

type Hub struct {

	Clients map[*Client]bool
	BroadcastMessagesCh chan []byte
	RegisterClientCh chan *Client
	UnregisterClientCh chan *Client
}


