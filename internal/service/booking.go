package service

import (
	"context"
	"fmt"
	"time"

	"applicationDesignTest/internal/logg"
	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/tools"
)

// BookingService is a service for booking rooms.
type PersistentRepo interface {
	GetAvailability(context.Context) ([]model.RoomAvailability, error)
	SaveOrder(context.Context, model.Order) error
	SetAvailability(context.Context, []model.RoomAvailability) error
	Lock() error
	UnLock() error
}

type BookingService struct {
	repo PersistentRepo
}

func NewBookingService(repo PersistentRepo) *BookingService {
	return &BookingService{
		repo: repo,
	}
}

func (s *BookingService) DoBookingOrder(ctx context.Context, order model.Order) error {
	s.repo.Lock()
	defer s.repo.UnLock()

	persistentAvailability, err := s.repo.GetAvailability(ctx)
	if err != nil {
		return fmt.Errorf("check availability error: %w", err)
	}

	err = doOrder(order, persistentAvailability)
	if err != nil {
		return fmt.Errorf("check order error: %w", err)
	}

	err = s.repo.SaveOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("save order error: %w", err)
	}

	err = s.repo.SetAvailability(ctx, persistentAvailability)
	if err != nil {
		return fmt.Errorf("set avalabilyty error: %w", err)
	}

	return nil
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
