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

	conn, err := upgrader.Upgrade(w, r, nil);
	if err != nil {

		log.Printf("Error while performing ws handshake: %v\n", err);
		return;
	}
	defer conn.Close()
	log.Println("Established websocket connetion");

	for {

		msg_typ, message, err := conn.ReadMessage()
		if err != nil {

			log.Printf("Error reading message in ws, closing conneciton: %v\n", err);
			break;
		}

		err_write := conn.WriteMessage(msg_typ, message);
		if err_write != nil {

			log.Printf("Failed writing message on websocket: %v\n", err_write);
			break
		}

	}

	return;
}
