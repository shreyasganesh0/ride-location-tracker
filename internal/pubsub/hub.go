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
	TopicMap 		    map[string]map[*Client]bool
	TopicSubCh			chan *Subscription
}

func NewHub() *Hub {

	var hub Hub;

	hub.Clients             = make(map[*Client]bool)
	hub.PublishMessagesCh   = make(chan *MessagePayload)
	hub.RegisterClientCh 	= make(chan *Client)
	hub.UnregisterClientCh  = make(chan *Client)
	hub.TopicMap		    = make(map[string]map[*Client]bool)
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

				for topic, clients := range h.TopicMap {

					if _, ok := clients[c_u]; ok {

						delete(clients, c_u)

						if len(clients) == 0 {

							delete(h.TopicMap, topic)
						}
					}
				}
				delete(h.Clients, c_u);
				close(c_u.OutboundMessagesCh)
			}

		case payload := <-h.PublishMessagesCh:

			if clients, ok := h.TopicMap[payload.Topic]; !ok {

				for client := range clients {

					client.OutboundMessagesCh <- &payload.Location
				}
			}

		case subs := <-h.TopicSubCh:

			 curr_topic := subs.Topic
			 c          := subs.Client;

			 if _, ok := h.TopicMap[subs.Topic]; !ok {
				 h.TopicMap[curr_topic] = make(map[*Client]bool)
			 }

			 h.TopicMap[curr_topic][c] = true;
			 log.Printf("Client %s subscribed to topic %s\n", c, curr_topic);

		}

	}
}


