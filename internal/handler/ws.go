package handler

import (
	"log"
	"net/http"
	"github.com/redis/go-redis/v9"
	"github.com/gorilla/websocket"
	"github.com/shreyasganesh0/ride-location-tracker/internal/broadcast"
)

var upgrader = websocket.Upgrader {

	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(req *http.Request) bool {

		return true //accept connections from anywhere for now
	},
};

func WsHandler(rdb *redis.Client, hub *broadcast.Hub, 
	w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil);
	if err != nil {

		log.Printf("Error while performing ws handshake: %v\n", err);
		return;
	}
	log.Println("Established websocket connetion");

	client := broadcast.NewClient(hub, conn)
	hub.RegisterClientCh<- client

	go client.ReadFromSocket(rdb)
	go client.WriteToSocket()

	return;
}
