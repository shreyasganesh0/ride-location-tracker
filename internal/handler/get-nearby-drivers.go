package handler

import (
	"log"
	"fmt"
	"context"
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/shreyasganesh0/ride-location-tracker/internal/pubsub"
)

func GetNearbyDriversHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	
	lon := r.FormValue("lon");
	if lon == "" {

		log.Println("Could not parse longitude from req");
		err_s := fmt.Sprintf("Invalid longitude value in request")
		err_code := 400
		WriteClientError(w, err_s, err_code);
		return;
	}
	flt_lon, err := strconv.ParseFloat(lon, 64);
	if err != nil {
		log.Println("Could not parse longitude from req to float: ", err);
		err_s := fmt.Sprintf("Invalid longitude value in request")
		err_code := 400
		WriteClientError(w, err_s, err_code);
		return;
	}

	lat := r.FormValue("lat");
	if lat == "" {

		log.Println("Could not parse latitude from req");
		err_s := fmt.Sprintf("Invalid latitude value in request")
		err_code := 400
		WriteClientError(w, err_s, err_code);
		return;

	}
	flt_lat, err := strconv.ParseFloat(lat, 64);
	if err != nil {
		log.Println("Could not parse latitude from req to float: ", err);
		err_s := fmt.Sprintf("Invalid latitude value in request")
		err_code := 400
		WriteClientError(w, err_s, err_code);
		return;

	}

	radius := r.FormValue("radius");
	if radius == "" {
		log.Println("Could not parse radius from req");
		err_s := fmt.Sprintf("Invalid radius value in request")
		err_code := 400
		WriteClientError(w, err_s, err_code);
		return;
	}
	flt_radius, err := strconv.ParseFloat(radius, 64);
	if err != nil {
		log.Println("Could not parse radius from req to float: ", err);
		err_s := fmt.Sprintf("Invalid radius value in request")
		err_code := 400
		WriteClientError(w, err_s, err_code);
		return;
	}


	res, err_q := rdb.GeoSearch(context.TODO(), "driver_locations",
		&redis.GeoSearchQuery{
			Longitude: flt_lon,
			Latitude: flt_lat,
			Radius: flt_radius,
			RadiusUnit: "km",
		}).Result();
	if err_q != nil {

		log.Println("Failed to get query for driver: ", err_q);
		err_s := fmt.Sprintf("Failed to get drivers");
		err_code := 400
		WriteClientError(w, err_s, err_code);
		return;
	}


	type MessageArr struct {

		Messages []pubsub.Message `json:"messages"`
	}

	msgs := MessageArr{
		Messages: make([]pubsub.Message,0),
	}

	for _, driverId := range res {
		
		message, err := GetDriverDetails(driverId, rdb, w, r)
		if err != nil {

			log.Println("Recieved error when getting driver details: \n", err);
			continue;
		}

		msgs.Messages = append(msgs.Messages, message)
	}

	msg_byts, err := json.Marshal(msgs)
	if err != nil {

		log.Println("Error parsing message to bytes: \n", err)
		err_s := fmt.Sprintf("Server failed to parse location\n");
		code := 500
		WriteClientError(w, err_s, code)
		return;
	}

	_, err_write := w.Write(msg_byts);
	if err_write != nil {

		log.Println("Error writing message: \n", err)
	}

	return;

}


