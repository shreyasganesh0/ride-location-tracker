package handler

import (
	"log"
	"os"
	"fmt"
	"net/http"
	"github.com/golang-jwt/jwt/v5"
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

	token_str := r.FormValue("token");
	if token_str == "" {

		log.Println("Couldnt get token from body")
		code := http.StatusUnauthorized
		err_s := fmt.Sprintf("Token not present\n")
		WriteClientError(w, err_s, code)
		return;
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(token_str, func(token *jwt.Token) (any, error) {

		return []byte(jwtSecret), nil

	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {

		log.Printf("Couldnt validate token: %v\n", err)
		code := http.StatusUnauthorized
		err_s := fmt.Sprintf("Token not authorized\n")
		WriteClientError(w, err_s, code)
		return;
	}

	claims, ok := token.Claims.(jwt.MapClaims); 
	if !ok {
		log.Printf("Couldnt validate token: %v\n", err)
		code := http.StatusUnauthorized
		err_s := fmt.Sprintf("Token not authorized\n")
		WriteClientError(w, err_s, code)
		return;
	}
	driverId := claims["did"].(string)

	conn, err := upgrader.Upgrade(w, r, nil);
	if err != nil {

		log.Printf("Error while performing ws handshake: %v\n", err);
		return;
	}
	log.Println("Established websocket connetion");

	client := broadcast.NewClient(driverId, hub, conn)
	hub.RegisterClientCh<- client

	go client.ReadFromSocket(rdb)
	go client.WriteToSocket()

	return;
}
