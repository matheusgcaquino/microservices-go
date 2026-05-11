package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/service"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer

	service *service.TripService
}

func NewGRPCHandler(server *grpc.Server, service *service.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()

	pickupCoord := &types.Coordinates{
		Latitude:  pickup.Latitude,
		Longitude: pickup.Longitude,
	}
	destinationCoord := &types.Coordinates{
		Latitude:  destination.Latitude,
		Longitude: destination.Longitude,
	}
	userID := req.GetUserID()

	rideFares, routes, err := h.service.PreviewTrip(ctx, userID, pickupCoord, destinationCoord)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return PreviewTripResponse(rideFares, routes), nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	trip, err := h.service.CreateTrip(ctx, req.GetUserID(), req.GetRideFareID())
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return CreateTripResponse(trip), nil
}
