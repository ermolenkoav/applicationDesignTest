package rest

import (
	"fmt"
	"net/http"

	"applicationDesignTest/internal/logg"
)

func (a *api) createOrder(w http.ResponseWriter, r *http.Request) {
	newOrder, err := parseOrder(r.Body)
	if err != nil {
		badRequest(w, "parse error")
		return
	}

	err = a.bService.DoBookingOrder(r.Context(), newOrder)
	if err != nil {
		badRequest(w, fmt.Sprintf("service error: %s", err))
		return
	}

	successRequest(w, newOrder)

	logg.Info("Order successfully created: %v", newOrder)
}
