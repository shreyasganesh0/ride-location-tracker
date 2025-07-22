package handler

import (
	"log"
	"net/http"
)

func WriteClientError(w http.ResponseWriter, resp string, code int) {

	w.WriteHeader(code)
	_, err := w.Write([]byte(resp));
	if err != nil {

		log.Printf("Error writing response bytes: %w", err);
	}
	return;
}
