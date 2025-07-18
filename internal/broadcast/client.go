package broadcast

import (
	"log"
	"fmt"
	"encoding/json"
	"context"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Client struct {

	Hub *Hub
	Conn *websocket.Conn
	OutboundMessagesCh chan *Message
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {

	log.Println("Creating a new client...")

	var client Client

	client.Hub = hub

	client.Conn = conn

	client.OutboundMessagesCh = make(chan *Message, 256)
	
	return &client

}

func (c *Client) CleanupClient() {

	c.Hub.UnregisterClientCh <- c  
	c.Conn.Close()
}

func (c *Client) ReadFromSocket(rdb *redis.Client) {

	defer c.CleanupClient()

	for {

		_, msg_byts, err := c.Conn.ReadMessage()
		if err != nil {

			log.Printf("Error reading message in ws, closing conneciton: %v\n", err);
			break;
		}

		var message Message
		err_json := json.Unmarshal(msg_byts, &message)
		if err_json != nil {

			log.Printf("Error marshalling a message %v\n", err);
			continue
		}


		key := fmt.Sprintf("driver:%s", message.DriverID);
		insert_location_map := map[string]interface{}{

			"longitude": message.Longitude,
			"latitude": message.Latitude,
		}
		err_set := rdb.HSet(context.Background(), key, insert_location_map).Err();
		if err_set != nil {

			log.Printf("Error uploding location of driver %s due to: %v\n",
				message.DriverID, err_set);
			continue;
		}

		c.Hub.BroadcastMessagesCh <- &message
	}
}

func (c *Client) WriteToSocket() {

	defer c.Conn.Close()

	for {
		select {

		case message, ok := <-c.OutboundMessagesCh:

			if !ok {

				c.Conn.WriteMessage(websocket.CloseMessage, []byte{});
				break
			}

			msg_byts, err_json := json.Marshal(message); //pointers encode to value?
			if err_json != nil {

				log.Printf("Recived unmarshallable message: %w\n", err_json);
				continue
			}

			err_write := c.Conn.WriteMessage(websocket.TextMessage, msg_byts);//maybe messagetype to be passed along with outbound message channel

			if err_write != nil {

				log.Printf("Failed writing message on websocket: %w\n", err_write);
				break
			}
		}
	}
}
