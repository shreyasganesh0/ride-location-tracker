package handler

import (
	"log"
	"net/http"
	"strconv"
	"github.com/redis/go-redis/v9"
	"github.com/shreyasganesh0/ride-location-tracker/internal/broadcast"
)

func GetNearbyDriversHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	
	lon := r.FormValue("lon");
	if lon == "" {

		err_s := log.Sprintf("Invalid longitude value in request")
		err_code := 400
		WriteClientError(&w, err_s, err_code);
		return;
	}
	flt_lon, err := strconv.ParseFloat(lon, 64);
	if err != nil {
		log.Printf("Could not parse longitude from req to float: %+v", err);
		err_s := log.Sprintf("Invalid longitude value in request")
		err_code := 400
		WriteClientError(&w, err_s, err_code);
		return;

	}

	lat := r.FormValue("lat");
	if lat == "" {

		err_s := log.Sprintf("Invalid latitude value in request")
		err_code := 400
		WriteClientError(&w, err_s, err_code);
		return;

	}
	flt_lat, err := strconv.ParseFloat(lat, 64);
	if err != nil {
		log.Printf("Could not parse longitude from req to float: %+v", err);
		err_s := log.Sprintf("Invalid longitude value in request")
		err_code := 400
		WriteClientError(&w, err_s, err_code);
		return;

	}

	radius := r.FormValue("radius");
	if radius == "" {
		err_s := log.Sprintf("Invalid radius value in request")
		err_code := 400
		WriteClientError(&w, err_s, err_code);
		return;
	}
	flt_radius, err := strconv.ParseFloat(radius, 64);
	if err != nil {
		log.Printf("Could not parse longitude from req to float: %+v", err);
		err_s := log.Sprintf("Invalid longitude value in request")
		err_code := 400
		WriteClientError(&w, err_s, err_code);
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

		log.Printf("Failed to get query for driver: %+v\n", err_q);
		err_s := log.Sprintf("Failed to get drivers");
		err_code := 400
		WriteClientError(&w, err_s, err_code);
		return;
	}


	for driverId := range res {
		
		driver_byts, err := GetDriverDetails(driverId, rdb, w, r)
		curr_msg := broadcast.Message{

			DriverID: driverId,

	}


}


