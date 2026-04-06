package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/shared/types"
)

type OsrmRouter struct {
}

type osrmResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func (router *OsrmRouter) GetRoutes(ctx context.Context, pickup, destination *types.Coordinates) (*types.Routes, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}

	var osrmResponse osrmResponse
	if err := json.Unmarshal(body, &osrmResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return router.routesResponse(osrmResponse), nil
}

func (router *OsrmRouter) routesResponse(osrmResponse osrmResponse) *types.Routes {
	var routesResponse types.Routes
	for _, route := range osrmResponse.Routes {
		newRoute := &types.Route{
			Distance: route.Distance,
			Duration: route.Duration,
			Geometry: &types.Geometry{},
		}

		for _, coord := range route.Geometry.Coordinates {
			newRoute.Geometry.Coordinates = append(
				newRoute.Geometry.Coordinates,
				&types.Coordinates{
					Longitude: coord[0],
					Latitude:  coord[1],
				},
			)
		}

		routesResponse.Routes = append(routesResponse.Routes, newRoute)
	}

	return &routesResponse
}
