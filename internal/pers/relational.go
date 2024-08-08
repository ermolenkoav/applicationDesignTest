package pers

import (
	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/tools"
)

type persistent struct {
	availability []model.RoomAvailability
	orders       []model.Order
}

func NewPersistent() *persistent {
	av := []model.RoomAvailability{
		{"reddison", "lux", tools.Date(2024, 1, 1), 1},
		{"reddison", "lux", tools.Date(2024, 1, 2), 1},
		{"reddison", "lux", tools.Date(2024, 1, 3), 1},
		{"reddison", "lux", tools.Date(2024, 1, 4), 1},
		{"reddison", "lux", tools.Date(2024, 1, 5), 0},
	}
	return &persistent{
		availability: av,
		orders:       make([]model.Order, 0),
	}
}

func (p *persistent) GetAvailability() ([]model.RoomAvailability, error) {
	return p.availability, nil
}

func (p *persistent) SaveOrder(newOrder model.Order) error {
	p.orders = append(p.orders, newOrder)
	return nil
}
