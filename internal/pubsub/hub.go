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
	TopicSubCh			chan *Subscription
}

func NewHub() *Hub {

	var hub Hub;

	hub.Clients             = make(map[*Client]bool)
	hub.PublishMessagesCh   = make(chan *MessagePayload)
	hub.RegisterClientCh 	= make(chan *Client)
	hub.UnregisterClientCh  = make(chan *Client)
	hub.TopicMap		    = make(map[string][]*Client)
	hub.TopicSubCh			= make(chan *Subscription);

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

			if clients, ok := h.TopicMap[payload.Topic]; !ok {

				for _, client := range clients {

					client.OutboundMessagesCh <- &payload.Location
				}
			}

		case subs := <-h.TopicSubCh:

			 curr_topic := subs.Topic
			 c          := subs.Client;

			 h.TopicMap[curr_topic] = append(c.Hub.TopicMap[curr_topic], c);

		}

	}
}


