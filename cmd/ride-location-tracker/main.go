package main

import (
	"log"
	"net/http"
	"github.com/shreyasganesh0/ride-location-tracker/internal/handler"
	"github.com/shreyasganesh0/ride-location-tracker/internal/broadcast"
	"github.com/shreyasganesh0/ride-location-tracker/internal/database"
)

func startup_message_handler() *broadcast.Hub{

	hub := broadcast.NewHub()
	go hub.Run()
	return hub
}

func main() {
	log.Println("Starting echo server...");

	hub := startup_message_handler()
	redis_client := database.NewRedisClient();

	http.HandleFunc("/" , handler.DefHandler);
	//http.HandleFunc("/echo" , handler.EchoHandler);
	http.HandleFunc("/ws" , func(w http.ResponseWriter, r *http.Request) { 
		handler.WsHandler(redis_client, hub, w, r)
	});
	http.HandleFunc("GET /api/drivers/{driverId}" , 
		func(w http.ResponseWriter, r *http.Request) {
			handler.GetDriverLocationHandler(redis_client, w, r)
	});
	log.Fatal(http.ListenAndServe(":8080", nil));
}


