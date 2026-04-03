package repository

import (
	"context"
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

func (repo *inmemRepository) SaveTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	repo.trips[trip.ID] = trip
	return trip, nil
}

func (repo *inmemRepository) SaveRideFare(ctx context.Context, rideFare *domain.RideFareModel) (*domain.RideFareModel, error) {
	repo.rideFares[rideFare.ID] = rideFare
	return rideFare, nil
}
