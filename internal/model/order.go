package model

import (
	"errors"
	"time"
)

type Order struct {
	HotelID   string    `json:"hotel_id"`
	RoomID    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

func (o Order) Validate() error {
	if o.HotelID == "" {
		return errors.New("hotel_id is required")
	}
	if o.RoomID == "" {
		return errors.New("room_id is required")
	}
	if o.UserEmail == "" {
		return errors.New("email is required")
	}
	if o.From.IsZero() {
		return errors.New("from date is required")
	}
	if o.To.IsZero() {
		return errors.New("to date is required")
	}
	if !o.From.Before(o.To) {
		return errors.New("from must be before to")
	}
	return nil
}
