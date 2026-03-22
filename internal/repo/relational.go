package repo

import (
	"context"

	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/tools"
)

// Repo is a thread-unsafe in-memory implementation of the booking repository.
// Callers are responsible for serializing concurrent access (e.g. via a service-level mutex).
type Repo struct {
	availability []model.RoomAvailability
	orders       []model.Order
}

func New() *Repo {
	return &Repo{
		availability: []model.RoomAvailability{
			{HotelID: "reddison", RoomID: "lux", Date: tools.Date(2024, 1, 1), Quota: 1},
			{HotelID: "reddison", RoomID: "lux", Date: tools.Date(2024, 1, 2), Quota: 1},
			{HotelID: "reddison", RoomID: "lux", Date: tools.Date(2024, 1, 3), Quota: 1},
			{HotelID: "reddison", RoomID: "lux", Date: tools.Date(2024, 1, 4), Quota: 1},
			{HotelID: "reddison", RoomID: "lux", Date: tools.Date(2024, 1, 5), Quota: 0},
		},
		orders: make([]model.Order, 0),
	}
}

func (r *Repo) GetAvailability(_ context.Context) ([]model.RoomAvailability, error) {
	result := make([]model.RoomAvailability, len(r.availability))
	copy(result, r.availability)
	return result, nil
}

func (r *Repo) SaveOrder(_ context.Context, order model.Order) error {
	r.orders = append(r.orders, order)
	return nil
}

func (r *Repo) SetAvailability(_ context.Context, data []model.RoomAvailability) error {
	r.availability = data
	return nil
}
