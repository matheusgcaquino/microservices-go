package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
)

type inmemRepository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (repo *inmemRepository) SaveTrip(ctx context.Context, trip *domain.TripModel) error {
	repo.trips[trip.ID] = trip
	return nil
}

func (repo *inmemRepository) SaveRideFare(ctx context.Context, rideFare *domain.RideFareModel) error {
	repo.rideFares[rideFare.ID] = rideFare
	return nil
}

func (repo *inmemRepository) GetRideFare(ctx context.Context, userID string, rideFareId string) (*domain.RideFareModel, error) {
	rideFare, exists := repo.rideFares[rideFareId]
	if !exists || rideFare.UserID != userID {
		return nil, fmt.Errorf("ride fare not found")
	}
	return rideFare, nil
}

func (r *inmemRepository) GetTripByID(ctx context.Context, id string) (*domain.TripModel, error) {
	trip, ok := r.trips[id]
	if !ok {
		return nil, nil
	}
	return trip, nil
}
