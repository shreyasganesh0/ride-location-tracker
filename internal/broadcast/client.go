package broadcast

import (
	"log"
	"github.com/gorilla/websocket"
)

type Client struct {

	Hub *Hub
	Conn *websocket.Conn
	OutboundMessagesCh chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {

	log.Println("Creating a new client...")

	var client Client

	client.Hub = hub

	client.Conn = conn

	client.OutboundMessagesCh = make(chan []byte, 256)

	return &client

}

func (c *Client) CleanupClient() {

	c.Hub.UnregisterClientCh <- c  
	c.Conn.Close()
}

func (c *Client) ReadFromSocket() {

	defer c.CleanupClient()

	for {

		_, message, err := c.Conn.ReadMessage()
		if err != nil {

			log.Printf("Error reading message in ws, closing conneciton: %v\n", err);
			break;
		}

		c.Hub.BroadcastMessagesCh <- message
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

			err_write := c.Conn.WriteMessage(websocket.TextMessage, message);//maybe messagetype to be passed along with outbound message channel
			if err_write != nil {

				log.Printf("Failed writing message on websocket: %w\n", err_write);
				break
			}
		}
	}
}
