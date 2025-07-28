package handler

import (
	"log"
	"fmt"
	"net/http"
	"strconv"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/shreyasganesh0/ride-location-tracker/internal/pubsub"
)

func GetDriverDetails(driverId string, rdb *redis.Client, w http.ResponseWriter, r *http.Request) (pubsub.Message, error) {

	var message pubsub.Message
	key := fmt.Sprintf("driver:%s", driverId);

	results, err := rdb.HGetAll(context.TODO(), key).Result()
	if err != nil {

		err_s := fmt.Errorf("Invalid driver request sent. Driver %s not sent\n", driverId);
		log.Printf("Error reading driver key %s: %w\n", key, err);
		code := 400
		WriteClientError(w, err_s.Error(), code)
		return message, err_s;
	}

	//message.DriverID = driverId

	for key, val := range results {

		switch (key) {

		case "longitude":
			
			longitude, err := strconv.ParseFloat(val, 64)
			if err != nil {

				log.Printf("Failed to convert longitude key to float: %w\n", err)
				err_s := fmt.Errorf("Failed to read location\n");
				code := 500
				WriteClientError(w, err_s.Error(), code)
				return message, err_s;
			}
			message.Longitude = longitude

		case "latitude":

			latitude, err := strconv.ParseFloat(val, 64)
			if err != nil {

				log.Printf("Failed to convert latitude key to float: %w\n", err)
				err_s := fmt.Errorf("Failed to read location\n");
				code := 500
				WriteClientError(w, err_s.Error(), code)
				return message, err_s;
			}
			message.Latitude = latitude
		default:

			log.Printf("Couldnt get appropriate keys:\t\nfound %s: %s, err: %w\n",
				key, val, err)
			err_s := fmt.Errorf("Driver details not found doesnt exist\n");
			code := 500
			WriteClientError(w, err_s.Error(), code)
			return message, err_s;
		}
	}

	return message, nil;
}

