package handler

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

func GetDriverLocationHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {

	driverId := r.PathValue("driverId");
	if driverId == "" {

		err_s := fmt.Sprintf("Error getting the driver id from path")
		log.Println(err_s)
		code := 400
		WriteClientError(w, err_s, code)
		return;
	}

	message, err := GetDriverDetails(driverId, rdb, w, r)
	if err != nil {

		log.Printf("Recieved error when getting driver details: %+v\n", err);
		return;
	}

	msg_byts, err := json.Marshal(message)
	if err != nil {

		err_s := fmt.Sprintf("Server failed to parse location\n");
		log.Printf("Error parsing message to bytes: %w\n", err)
		code := 500
		WriteClientError(w, err_s, code)
		return;
	}
	_, err_write := w.Write(msg_byts);
	if err_write != nil {

		log.Printf("Error writing message: %w\n", err)
	}

	return;

}
