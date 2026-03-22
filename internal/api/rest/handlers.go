package rest

import (
	"fmt"
	"net/http"

	"applicationDesignTest/internal/logg"
)

// @Summary		Booking
// @Description	Create a hotel room booking order
// @Router			/api/v1/orders	[post]
// @Param			request			body	model.Order	true	"hotel_id, room_id, email, from, to"
// @Produce		json
// @Success		201	{object}	model.Order
// @Failure		400	{object}	string
// @Failure		422	{object}	string
func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	newOrder, err := parseOrder(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err = s.bService.DoBookingOrder(r.Context(), newOrder); err != nil {
		writeError(w, http.StatusUnprocessableEntity, fmt.Sprintf("booking failed: %s", err))
		return
	}

	writeJSON(w, http.StatusCreated, newOrder)
	logg.Info("order created: %v", newOrder)
}
