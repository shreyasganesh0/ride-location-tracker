package handler

import (
	"log"
	"fmt"
	"net/http"
	"strconv"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/shreyasganesh0/ride-location-tracker/internal/broadcast"
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

	key := fmt.Sprintf("driver:%s", driverId);

	results, err := rdb.HGetAll(context.Background(), key).Result()
	if err != nil {

		err_s := fmt.Sprintf("Invalid driver request sent. Driver %s not sent\n", driverId);
		log.Printf("Error reading driver key %s: %w\n", key, err);
		code := 400
		WriteClientError(w, err_s, code)
		return;
	}

	var message broadcast.Message
	message.DriverID = driverId

	for key, val := range results {

		switch (key) {

		case "longitude":
			
			longitude, err := strconv.ParseFloat(val, 64)
			if err != nil {

				log.Printf("Failed to convert longitude key to float: %w\n", err)
				err_s := fmt.Sprintf("Failed to read location\n");
				code := 500
				WriteClientError(w, err_s, code)
				return;
			}
			message.Longitude = longitude

		case "latitude":

			latitude, err := strconv.ParseFloat(val, 64)
			if err != nil {

				log.Printf("Failed to convert latitude key to float: %w\n", err)
				err_s := fmt.Sprintf("Failed to read location\n");
				code := 500
				WriteClientError(w, err_s, code)
				return;
			}
			message.Latitude = latitude
		default:

			err_s := fmt.Sprintf("Driver details not found doesnt exist\n");
			log.Printf("Couldnt get appropriate keys:\t\nfound %s: %s, err: %w\n",
				key, val, err)
			code := 500
			WriteClientError(w, err_s, code)
			return;
		}
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
