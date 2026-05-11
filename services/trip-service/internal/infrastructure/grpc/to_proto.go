package grpc

import (
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

func RideFareToProto(r *domain.RideFareModel) *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID,
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func ToRideFaresProto(ridefares []*domain.RideFareModel) []*pb.RideFare {
	var protoFares []*pb.RideFare
	for _, fare := range ridefares {
		protoFares = append(protoFares, RideFareToProto(fare))
	}
	return protoFares
}

func RoutesToProto(o *types.Routes) *pb.Route {
	route := o.Routes[0]
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
