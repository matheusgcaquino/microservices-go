package domain

import "context"

type RideFareModel struct {
	ID                string
	UserID            string
	PackafeSlug       string
	TotalPriceInCents float64
}

func NewRideFare(
	ctx context.Context,
	id string,
	userID string,
	packafeSlug string,
	totalPriceInCents float64) (*RideFareModel, error) {
	return &RideFareModel{
		ID:                id,
		UserID:            userID,
		PackafeSlug:       packafeSlug,
		TotalPriceInCents: totalPriceInCents,
	}, nil
}