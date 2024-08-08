package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"applicationDesignTest/internal/logg"
	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/tools"
)

func buildOrder(r io.Reader) (model.Order, error) {
	var newOrder model.Order
	err := json.NewDecoder(r).Decode(&newOrder)
	if err != nil {
		return newOrder, err
	}
	return newOrder, nil
}

func doOrder(newOrder model.Order, persistentAvailability []model.RoomAvailability) error {
	daysToBook := tools.DaysBetween(newOrder.From, newOrder.To)

	unavailableDays := make(map[time.Time]struct{})
	for _, day := range daysToBook {
		unavailableDays[day] = struct{}{}
	}

	for _, dayToBook := range daysToBook {
		for i, availability := range persistentAvailability {
			if !availability.Date.Equal(dayToBook) || availability.Quota < 1 {
				continue
			}
			availability.Quota -= 1
			persistentAvailability[i] = availability
			delete(unavailableDays, dayToBook)
		}
	}

	if len(unavailableDays) != 0 {
		errMessage := fmt.Sprintf("Hotel room is not available for selected dates:\n%v\n%v", newOrder, unavailableDays)
		logg.Errorf(errMessage)
		return fmt.Errorf(errMessage)
	}

	return nil
}

func (a *api) createOrder(w http.ResponseWriter, r *http.Request) {
	newOrder, err := buildOrder(r.Body)
	if err != nil {
		badRequest(w, "build error")
		return
	}

	persistentAvailability, err := a.repo.GetAvailability()
	if err != nil {
		badRequest(w, "check availability error")
		return
	}

	err = doOrder(newOrder, persistentAvailability)
	if err != nil {
		badRequest(w, "check order error")
		return
	}

	err = a.repo.SaveOrder(newOrder)
	if err != nil {
		badRequest(w, "save order error")
		return
	}

	successRequest(w, newOrder)

	logg.Info("Order successfully created: %v", newOrder)
}
