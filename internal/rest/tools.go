package rest

import (
	"encoding/json"
	"net/http"
)

func badRequest(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}

func successRequest(w http.ResponseWriter, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}
