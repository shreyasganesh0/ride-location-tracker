package main

import (
	"log"
	"net/http"
	"github.com/shreyasganesh0/ride-location-tracker/internal/handler"
)


func main() {
	log.Println("Starting echo server...");

	http.HandleFunc("/" , handler.DefHandler);
	http.HandleFunc("/echo" , handler.EchoHandler);
	http.HandleFunc("/ws" , handler.WsHandler);
	log.Fatal(http.ListenAndServe(":8080", nil));
}


