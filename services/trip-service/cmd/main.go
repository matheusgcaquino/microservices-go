package main

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	inmemRepository := repository.NewInmemRepository()

	rideFareService := service.NewRideFareService(inmemRepository)
	rideFare, _ := rideFareService.CreateRideFare(ctx, uuid.NewString(), uuid.NewString(), 234.656)

	log.Print("Ride: ")
	log.Println(rideFare)

	tripService := service.NewTripService(inmemRepository)
	trip, _ := tripService.CreateTrip(ctx, rideFare)

	log.Print("Trip: ")
	log.Println(trip)

	for {
		time.Sleep(time.Second)
	}
}
