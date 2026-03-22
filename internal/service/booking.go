package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/tools"
)

type BookingRepo interface {
	GetAvailability(context.Context) ([]model.RoomAvailability, error)
	SaveOrder(context.Context, model.Order) error
	SetAvailability(context.Context, []model.RoomAvailability) error
}

type BookingService struct {
	repo BookingRepo
	mu   sync.Mutex
}

func NewBookingService(repo BookingRepo) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) DoBookingOrder(ctx context.Context, order model.Order) error {
	if err := order.Validate(); err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	availability, err := s.repo.GetAvailability(ctx)
	if err != nil {
		return fmt.Errorf("get availability: %w", err)
	}

	if err = reserveDays(order, availability); err != nil {
		return err
	}

	if err = s.repo.SaveOrder(ctx, order); err != nil {
		return fmt.Errorf("save order: %w", err)
	}

	if err = s.repo.SetAvailability(ctx, availability); err != nil {
		return fmt.Errorf("set availability: %w", err)
	}

	return nil
}

// reserveDays decrements quota for each day of the booking.
// Returns an error if any required day has no available quota.
func reserveDays(order model.Order, availability []model.RoomAvailability) error {
	daysToBook := tools.DaysBetween(order.From, order.To)

	unavailable := make(map[time.Time]struct{}, len(daysToBook))
	for _, d := range daysToBook {
		unavailable[d] = struct{}{}
	}

	for _, day := range daysToBook {
		for i := range availability {
			a := &availability[i]
			if a.HotelID != order.HotelID || a.RoomID != order.RoomID {
				continue
			}
			if !a.Date.Equal(day) || a.Quota < 1 {
				continue
			}
			a.Quota--
			delete(unavailable, day)
			break
		}
	}

	if len(unavailable) != 0 {
		return fmt.Errorf("room %s/%s not available for dates: %v", order.HotelID, order.RoomID, unavailable)
	}

	return nil
}
