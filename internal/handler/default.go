package handler

import (
	"net/http"
)

func DefHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
