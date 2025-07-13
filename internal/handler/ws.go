package handler

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {

	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(req *http.Request) bool {

		return true //accept connections from anywhere for now
	},
};

func WsHandler(w http.ResponseWriter, r *http.Request) {

	_, err := upgrader.Upgrade(w, r, nil);
	if err != nil {

		log.Println("Error while performing ws handshake:", err);
		return;
	}

	log.Println("Established websocket connetion");
	return;
}
