package handler

import (
	"os"
	"fmt"
	"log"
	"time"
	"net/http"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
)


func AuthDriverHandler(w http.ResponseWriter, r *http.Request) {
	
	hmacSecret := os.Getenv("JWT_SECRET")
	if hmacSecret == "" {

		log.Println("Couldnt load jwt secret");
		code := 500
		err_s := fmt.Sprintf("Server couldnt authenticate\n");
		WriteClientError(w, err_s, code);
		return;
	}

	driverId := r.PostFormValue("driverId")
	if driverId == ""{
		
		log.Println("Couldnt get driverId secret");
		code := 400
		err_s := fmt.Sprintf("Invalid message body\n");
		WriteClientError(w, err_s, code);
		return;
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{

		"exp": jwt.NewNumericDate(time.Now().Add(1 * time.Hour)) ,
		"iss": "ride-location-tracker",
		"did": driverId,
	});

	ss, err := token.SignedString([]byte(hmacSecret));
	
	type JwtBody struct {
		Token string `json:"token"`
	}

	var jwt_resp = JwtBody{
		Token: ss,
	}

	jwt_resp_byts, err := json.Marshal(jwt_resp);
	if err != nil {

		log.Printf("Error creating response: %+v\n", err);
		code := 500
		err_s := fmt.Sprintf("Server error\n");
		WriteClientError(w, err_s, code);
		return
	}


	_, err = w.Write(jwt_resp_byts)
	if err != nil {

		log.Printf("Error creating response: %+v\n", err);
	}
}
