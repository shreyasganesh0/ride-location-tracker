package handler

import (
	"log"
	"net/http"
)

func DefHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("This is the default response."));
	if err != nil {
		log.Fatal(err.Error());
	}
}
