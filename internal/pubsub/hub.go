package pubsub

import (
	"log"
	//"github.com/gorilla/websocket"
)

type Hub struct {

	Clients				map[*Client]bool
	PublishMessagesCh   chan *MessagePayload
	RegisterClientCh    chan *Client
	UnregisterClientCh  chan *Client
	TopicMap 		    map[string][]*Client
}

func NewHub() *Hub {

	var hub Hub;

	hub.Clients             = make(map[*Client]bool)
	hub.PublishMessagesCh   = make(chan *MessagePayload)
	hub.RegisterClientCh 	= make(chan *Client)
	hub.UnregisterClientCh  = make(chan *Client)
	hub.TopicMap		    = make(map[string][]*Client)

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

		case payload := <-h.PublishMessagesCh:

			for _, c := range h.TopicMap[payload.Topic] {

				c.OutboundMessagesCh <- &payload.Location
			}
		}

	}
}


