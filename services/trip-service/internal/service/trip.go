package service

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"

	"github.com/google/uuid"
)

type TripRepository interface {
	SaveTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error)
}

type Router interface {
	GetRoutes(ctx context.Context, pickup, destination *types.Coordinates) (*types.Routes, error)
}

type TripService struct {
	repo   TripRepository
	router Router
}

func NewTripService(repo TripRepository, router Router) *TripService {
	return &TripService{
		repo:   repo,
		router: router,
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

func (s *TripService) GetRoutes(ctx context.Context, pickup, destination *types.Coordinates) (*types.Routes, error) {
	return s.router.GetRoutes(ctx, pickup, destination)
}
