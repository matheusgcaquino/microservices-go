package service

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"

	"github.com/google/uuid"
)

type RideFareRepository interface {
	SaveRideFare(ctx context.Context, rideFare *domain.RideFareModel) (*domain.RideFareModel, error)
}

type RideFareService struct {
	repo RideFareRepository
}

func NewRideFareService(repo RideFareRepository) *RideFareService {
	return &RideFareService{
		repo: repo,
	}
}

func (s *RideFareService) CreateRideFare(
	ctx context.Context,
	userID string,
	packafeSlug string,
	totalPriceInCents float64) (*domain.RideFareModel, error) {

	rideFare, error := domain.NewRideFare(ctx, uuid.NewString(), userID, packafeSlug, totalPriceInCents)
	if error != nil {
		return nil, error
	}

	s.repo.SaveRideFare(ctx, rideFare)
	return rideFare, nil
}
