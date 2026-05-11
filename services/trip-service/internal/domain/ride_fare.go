package domain

import (
	"context"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                string
	UserID            string
	PackageSlug       string
	TotalPriceInCents float64
}

func NewRideFare(
	ctx context.Context,
	id string,
	userID string,
	PackageSlug string,
	totalPriceInCents float64) (*RideFareModel, error) {
	return &RideFareModel{
		ID:                id,
		UserID:            userID,
		PackageSlug:       PackageSlug,
		TotalPriceInCents: totalPriceInCents,
	}, nil
}

func GetRideFares(ctx context.Context, userID string, routes *types.Routes) []*RideFareModel {
	baseFares := getBaseFares()
	rideFares := make([]*RideFareModel, len(baseFares))

	for i, fare := range baseFares {
		estimatedFare := calculateFare(fare, routes)
		rideFares[i], _ = NewRideFare(
			ctx,
			primitive.NewObjectID().Hex(),
			userID,
			estimatedFare.PackageSlug,
			estimatedFare.TotalPriceInCents,
		)
	}

	return rideFares
}

func calculateFare(fare *tripTypes.RideFares, routes *types.Routes) *tripTypes.RideFares {
	pricingCfg := tripTypes.DefaultPricingConfig()
	carPackagePrice := fare.TotalPriceInCents

	distanceKm := routes.Routes[0].Distance
	durationInMinutes := routes.Routes[0].Duration

	distanceFare := distanceKm * pricingCfg.PricePerUnitOfDistance
	timeFare := durationInMinutes * pricingCfg.PricingPerMinute
	totalPrice := carPackagePrice + distanceFare + timeFare

	return &tripTypes.RideFares{
		TotalPriceInCents: totalPrice,
		PackageSlug:       fare.PackageSlug,
	}
}

func getBaseFares() []*tripTypes.RideFares {
	return []*tripTypes.RideFares{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 1500,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 550,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 1000,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 2000,
		},
	}
}
