package broadcast

import (
	"log"
	//"github.com/gorilla/websocket"
)

type Hub struct {

	Clients map[*Client]bool
	BroadcastMessagesCh chan []byte
	RegisterClientCh chan *Client
	UnregisterClientCh chan *Client
}

func NewHub() *Hub {

	var hub Hub;

	hub.Clients = make(map[*Client]bool)
	hub.BroadcastMessagesCh = make(chan []byte)
	hub.RegisterClientCh = make(chan *Client)
	hub.UnregisterClientCh = make(chan *Client)

	return &hub
}

func (h *Hub) Run() {

	log.Println("Running the hub...")
	for {

		select {

		case c_r := <-h.RegisterClientCh:

			h.Clients[c_r] = true;

		case c_u := <-h.UnregisterClientCh:

			_, exists := h.Clients[c_u]
			if exists {

				delete(h.Clients, c_u);
				close(c_u.OutboundMessagesCh)
			}

		case msg := <-h.BroadcastMessagesCh:

			for client, _ := range h.Clients {
				
				select {
				case client.OutboundMessagesCh <- msg: //maybe need to pre allocate size
				default :
					// assujme client is dead
					close(client.OutboundMessagesCh)
					delete(h.Clients, client)
				}
			}
		}

	}
}


