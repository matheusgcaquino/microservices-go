package service

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"

	"github.com/google/uuid"
)

type TripRepository interface {
	SaveTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error)
}

type TripService struct {
	repo TripRepository
}

func NewTripService(repo TripRepository) *TripService {
	return &TripService{
		repo: repo,
	}
}

func (s *TripService) CreateTrip(ctx context.Context, rideFare *domain.RideFareModel) (*domain.TripModel, error) {

	trip, error := domain.NewTrip(ctx, uuid.NewString(), rideFare)
	if error != nil {
		return nil, error
	}

	s.repo.SaveTrip(ctx, trip)
	return trip, nil
}
