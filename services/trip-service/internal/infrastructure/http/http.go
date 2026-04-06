package http

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/types"
)

type HttpHandler struct {
	Service service.TripService
}

type previewTripRequest struct {
	UserID      string      `json:"userID"`
	Pickup      Coordinates `json:"pickup"`
	Destination Coordinates `json:"destination"`
}

type Coordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (s *HttpHandler) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	pickup := types.Coordinates{
		Longitude: reqBody.Pickup.Longitude,
		Latitude:  reqBody.Pickup.Latitude,
	}

	destination := types.Coordinates{
		Longitude: reqBody.Destination.Longitude,
		Latitude:  reqBody.Destination.Latitude,
	}
	t, err := s.Service.GetRoutes(ctx, &pickup, &destination)
	if err != nil {
		log.Println(err)
	}

	writeJSON(w, http.StatusOK, t)
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
