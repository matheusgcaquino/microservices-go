package events

import (
	"context"
	"encoding/json"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	pbd "ride-sharing/shared/proto/driver"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"github.com/rabbitmq/amqp091-go"
)

type driverConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  *service.TripService
}

func NewDriverConsumer(rabbitmq *messaging.RabbitMQ, service *service.TripService) *driverConsumer {
	return &driverConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (c *driverConsumer) Listen() error {
	return c.rabbitmq.ConsumeMessages(messaging.DriverTripResponseQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		var payload messaging.DriverTripResponseData
		if err := json.Unmarshal(message.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		log.Printf("driver response received message: %+v", payload)

		switch msg.RoutingKey {
		case contracts.DriverCmdTripAccept:
			if err := c.handleTripAccepted(ctx, payload.TripID, payload.Driver); err != nil {
				log.Printf("Failed to handle the trip accept: %v", err)
				return err
			}
		case contracts.DriverCmdTripDecline:
			if err := c.handleTripDeclined(ctx, payload.TripID, payload.RiderID); err != nil {
				log.Printf("Failed to handle the trip decline: %v", err)
				return err
			}
			return nil
		}
		log.Printf("unknown trip event: %+v", payload)

		return nil
	})
}

func (c *driverConsumer) handleTripAccepted(ctx context.Context, tripID string, driver *pbd.Driver) error {
	driverModel := &domain.DriverModel{
		ID:             driver.Id,
		Name:           driver.Name,
		ProfilePicture: driver.ProfilePicture,
		CarPlate:       driver.CarPlate,
	}

	trip, err := c.service.TripAccepted(ctx, tripID, driverModel)
	if err != nil {
		log.Printf("Failed to accept the trip: %v", err)
		return err
	}

	// Driver has been assigned -> publish this event to RB
	marshalledTrip, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	// Notify the rider that a driver has been assigned
	if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventDriverAssigned, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    marshalledTrip,
	}); err != nil {
		return err
	}

	// TODO: Notify the payment service to start a payment link

	return nil
}

func (c *driverConsumer) handleTripDeclined(ctx context.Context, tripID, riderID string) error {
	// When a driver declines, we should try to find another driver

	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	newPayload := messaging.TripEventData{
		Trip: TripToProto(trip),
	}

	marshalledPayload, err := json.Marshal(newPayload)
	if err != nil {
		return err
	}

	if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventDriverNotInterested,
		contracts.AmqpMessage{
			OwnerID: riderID,
			Data:    marshalledPayload,
		},
	); err != nil {
		return err
	}

	return nil
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

func TripToProto(trip *domain.TripModel) *pb.Trip {
	var protoDriver *pb.TripDriver
	if trip.Driver != nil {
		protoDriver = DriverToProto(trip.Driver)
	}

	return &pb.Trip{
		Id:           trip.ID,
		UserID:       trip.UserID,
		SelectedFare: RideFareToProto(trip.RideFare),
		Status:       trip.Status,
		Driver:       protoDriver,
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
