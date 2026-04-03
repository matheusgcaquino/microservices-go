package domain

import "context"

type TripModel struct {
	ID       string
	UserID   string
	Status   string
	RideFare *RideFareModel
}

func NewTrip(
	ctx context.Context,
	id string,
	rideFare *RideFareModel) (*TripModel, error) {
	return &TripModel{
		ID:       id,
		UserID:   rideFare.UserID,
		Status:   "PENDING",
		RideFare: rideFare,
	}, nil
}
