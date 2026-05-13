package grpc

import (
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

func PreviewTripResponse(rideFares []*domain.RideFareModel, routes *types.Routes) *pb.PreviewTripResponse {
	protoRoute := RouteToProto(routes.Routes[0])

	var protoFares []*pb.RideFare
	for _, fare := range rideFares {
		protoFares = append(protoFares, RideFareToProto(fare))
	}

	return &pb.PreviewTripResponse{
		Route:     protoRoute,
		RideFares: protoFares,
	}
}

func RouteToProto(route *types.Route) *pb.Route {
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
		protoDriver = DriverToProto(trip.Driver)
	}

	return &pb.CreateTripResponse{
		TripID: trip.ID,
		Trip: &pb.Trip{
			Id:           trip.ID,
			UserID:       trip.UserID,
			Status:       trip.Status,
			SelectedFare: RideFareToProto(trip.RideFare),
			Driver:       protoDriver,
		},
	}
}

func RideFareToProto(fare *domain.RideFareModel) *pb.RideFare {
	return &pb.RideFare{
		Id:                fare.ID,
		UserID:            fare.UserID,
		PackageSlug:       fare.PackageSlug,
		TotalPriceInCents: fare.TotalPriceInCents,
		Route:             RouteToProto(fare.Route),
	}
}

func DriverToProto(driver *domain.DriverModel) *pb.TripDriver {
	return &pb.TripDriver{
		Id:             driver.ID,
		Name:           driver.Name,
		ProfilePicture: driver.ProfilePicture,
		CarPlate:       driver.CarPlate,
	}
}

func TripToProto(t *domain.TripModel) *pb.Trip {
	return &pb.Trip{
		Id:           t.ID,
		UserID:       t.UserID,
		SelectedFare: RideFareToProto(t.RideFare),
		Status:       t.Status,
		Driver:       DriverToProto(t.Driver),
	}
}
