package handler

import (
	"net/http"
)

func HandleChat(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
