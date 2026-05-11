package domain

import "context"

type DriverModel struct {
	ID             string
	Name           string
	ProfilePicture string
	CarPlate       string
}

func NewDriver(
	ctx context.Context,
	id string,
	name string,
	profilePicture string,
	carPlate string) (*DriverModel, error) {
	return &DriverModel{
		ID:             id,
		Name:           name,
		ProfilePicture: profilePicture,
		CarPlate:       carPlate,
	}, nil
}
