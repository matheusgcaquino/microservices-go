package service

import (
	"context"
	"fmt"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"

	"github.com/google/uuid"
)

type TripRepository interface {
	SaveTrip(ctx context.Context, trip *domain.TripModel) error
	SaveRideFare(ctx context.Context, rideFare *domain.RideFareModel) error
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

func (s *TripService) GetRoutes(ctx context.Context, pickup *types.Coordinates, destination *types.Coordinates) (*types.Routes, error) {
	return s.router.GetRoutes(ctx, pickup, destination)
}

func (s *TripService) PreviewTrip(ctx context.Context, userID string, pickup *types.Coordinates, destination *types.Coordinates) ([]*domain.RideFareModel, *types.Routes, error) {
	routes, err := s.GetRoutes(ctx, pickup, destination)
	if err != nil {
		log.Println(err)
		return nil, nil, fmt.Errorf("failed to get route: %v", err)
	}

	rideFares := domain.GetRideFares(ctx, userID, routes)
	if len(rideFares) == 0 {
		return nil, nil, fmt.Errorf("failed to calculate fare")
	}

	for _, fare := range rideFares {
		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, nil, fmt.Errorf("failed to save trip fare: %w", err)
		}
	}

	return rideFares, routes, nil
}
