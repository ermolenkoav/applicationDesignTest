package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"applicationDesignTest/internal/model"
)

func parseOrder(r io.Reader) (model.Order, error) {
	var order model.Order
	if err := json.NewDecoder(r).Decode(&order); err != nil {
		return order, err
	}
	return order, nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
