package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"applicationDesignTest/internal/model"
)

func parseOrder(r io.Reader) (model.Order, error) {
	var newOrder model.Order
	err := json.NewDecoder(r).Decode(&newOrder)
	if err != nil {
		return newOrder, err
	}
	return newOrder, nil
}

func badRequest(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
}

func successRequest(w http.ResponseWriter, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}
