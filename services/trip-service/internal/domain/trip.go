package domain

import "context"

type TripModel struct {
	ID       string
	UserID   string
	Status   string
	RideFare *RideFareModel
	Driver   *DriverModel
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
		Driver:   nil,
	}, nil
}

func (t *TripModel) NewDriver(ctx context.Context, driver *DriverModel) error {
	t.Status = "ACCEPTED"
	t.Driver = driver
	return nil
}
