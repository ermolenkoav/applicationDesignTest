package repo

import (
	"context"
	"errors"
	"sync"

	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/tools"
)

type persistent struct {
	availability []model.RoomAvailability
	orders       []model.Order
	mt           *sync.RWMutex
	isOpen       bool
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
		mt:           &sync.RWMutex{},
	}
}

func (p *persistent) GetAvailability(_ context.Context) ([]model.RoomAvailability, error) {
	if p.isOpen {
		return p.availability, nil
	}
	return nil, errors.New("busy")
}

func (p *persistent) SaveOrder(_ context.Context, newOrder model.Order) error {
	if p.isOpen {
		p.orders = append(p.orders, newOrder)
		return nil
	}
	return errors.New("busy")
}

func (p *persistent) SetAvailability(_ context.Context, data []model.RoomAvailability) error {
	if p.isOpen {
		p.availability = data
		return nil
	}
	return errors.New("busy")
}

func (p *persistent) Lock() error {
	p.mt.Lock()
	p.isOpen = true
	return nil
}

func (p *persistent) UnLock() error {
	p.mt.Unlock()
	p.isOpen = false
	return nil
}
