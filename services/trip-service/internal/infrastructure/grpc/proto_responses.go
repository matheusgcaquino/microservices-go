package grpc

import (
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

func PreviewTripResponse(rideFares []*domain.RideFareModel, routes *types.Routes) *pb.PreviewTripResponse {
	protoRoute := routeToProto(routes.Routes[0])

	var protoFares []*pb.RideFare
	for _, fare := range rideFares {
		protoFares = append(protoFares, &pb.RideFare{
			Id:                fare.ID,
			UserID:            fare.UserID,
			PackageSlug:       fare.PackageSlug,
			TotalPriceInCents: fare.TotalPriceInCents,
			Route:             protoRoute,
		})
	}

	return &pb.PreviewTripResponse{
		Route:     protoRoute,
		RideFares: protoFares,
	}
}

func routeToProto(route *types.Route) *pb.Route {
	geometry := route.Geometry.Coordinates
	coordinates := make([]*pb.Coordinate, len(geometry))
	for i, coord := range geometry {
		coordinates[i] = &pb.Coordinate{
			Latitude:  coord.Latitude,
			Longitude: coord.Longitude,
		}
	}
	return &pb.Route{
		Geometry: []*pb.Geometry{
			{
				Coordinates: coordinates,
			},
		},
		Distance: route.Distance,
		Duration: route.Duration,
	}
}

func CreateTripResponse(trip *domain.TripModel) *pb.CreateTripResponse {
	var protoDriver *pb.TripDriver
	if trip.Driver != nil {
		protoDriver = &pb.TripDriver{
			Id:             trip.Driver.ID,
			Name:           trip.Driver.Name,
			ProfilePicture: trip.Driver.ProfilePicture,
			CarPlate:       trip.Driver.CarPlate,
		}
	}

	return &pb.CreateTripResponse{
		TripID: trip.ID,
		Trip: &pb.Trip{
			Id:     trip.ID,
			UserID: trip.UserID,
			Status: trip.Status,
			SelectedFare: &pb.RideFare{
				Id:                trip.RideFare.ID,
				UserID:            trip.RideFare.UserID,
				PackageSlug:       trip.RideFare.PackageSlug,
				TotalPriceInCents: trip.RideFare.TotalPriceInCents,
				Route:             routeToProto(trip.RideFare.Route),
			},
			Driver: protoDriver,
		},
	}
}
