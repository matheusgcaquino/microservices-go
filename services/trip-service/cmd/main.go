package main

import (
	"log"
	"net/http"
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
)

func main() {
	inmemRepository := repository.NewInmemRepository()
	tripService := service.NewTripService(inmemRepository, &h.OsrmRouter{})
	mux := http.NewServeMux()

	httphandler := h.HttpHandler{Service: *tripService}
	mux.HandleFunc("POST /preview", httphandler.HandleTripPreview)
	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server error: %v", err)

	}
}
