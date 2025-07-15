package main

import (
	"log"
	"net/http"
	"github.com/shreyasganesh0/ride-location-tracker/internal/handler"
	"github.com/shreyasganesh0/ride-location-tracker/internal/broadcast"
)

func startup_message_handler() *broadcast.Hub{

	hub := broadcast.NewHub()
	go hub.Run()
	return hub
}

func main() {
	log.Println("Starting echo server...");

	hub := startup_message_handler()

	http.HandleFunc("/" , handler.DefHandler);
	//http.HandleFunc("/echo" , handler.EchoHandler);
	http.HandleFunc("/ws" , func(w http.ResponseWriter, r *http.Request) { 
		handler.WsHandler(hub, w, r)
	});
	log.Fatal(http.ListenAndServe(":8080", nil));
}


