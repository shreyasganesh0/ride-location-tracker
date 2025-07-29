package pubsub

import (
	"log"
	"fmt"
	"encoding/json"
	"context"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/mmcloughlin/geohash"
)


type Client struct {

	Hub *Hub
	Conn *websocket.Conn
	OutboundMessagesCh chan *LatLng
	DriverID string
}

func NewClient(driverId string, hub *Hub, conn *websocket.Conn) *Client {

	log.Println("Creating a new client...")

	var client Client

	client.Hub = hub

	client.Conn = conn

	client.DriverID = driverId

	client.OutboundMessagesCh = make(chan *LatLng, 256)

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

		geo_key := "driver_locations" 

		err_add := rdb.GeoAdd(context.TODO(), geo_key, &redis.GeoLocation{

			Longitude: message.Payload.Location.Longitude,
			Latitude: message.Payload.Location.Latitude,
			Name: c.DriverID,
		}).Err();
		if err_add != nil {

			log.Printf("Error uploding location of driver %s due to: %v\n",
				c.DriverID, err_add);
			continue;
		}

		driver_insert_data := map[string]interface{}{
			"longitude": message.Payload.Location.Longitude,
			"latitude": message.Payload.Location.Latitude,
		}


		hset_key := fmt.Sprintf("driver:%s", c.DriverID);
		err_set := rdb.HSet(context.TODO(), hset_key, driver_insert_data).Err();
		if err_set != nil {

			log.Printf("Error uploding location of driver %s due to: %v\n",
				c.DriverID, err_set);
			continue;
		}

		switch (message.Type) {
		case Subscribe:

			curr_topic := message.Payload.Topic
			c.Hub.TopicMap[curr_topic] = append(c.Hub.TopicMap[curr_topic], c);

		case PublishLocation:

			topic := geohash.EncodeWithPrecision(message.Payload.Location.Latitude, 
				message.Payload.Location.Longitude, 6);

			message.Payload.Topic = topic
			c.Hub.PublishMessagesCh <- &message.Payload
		}

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

			err_write := c.Conn.WriteMessage(websocket.TextMessage, msg_byts);

			if err_write != nil {

				log.Printf("Failed writing message on websocket: %w\n", err_write);
				break
			}
		}
	}
}
