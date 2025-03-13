package maps

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/env"
	routespb "cloud.google.com/go/maps/routing/apiv2/routingpb"
	"context"
	"crypto/tls"
	"google.golang.org/api/option"
	"google.golang.org/api/places/v1"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

type MapsIf interface {
	GetAutocomplete(req dto.LocationRequest) ([]dto.LocationResponse, error)
	GetPlaceDetails(placeId string) (dto.PlaceDetails, error)
	GetDistance(from dto.Location, to dto.Location) (*int32, error)
	GenerateGoogleMapsURL(placeId string) string
}

type Maps struct {
	env           *env.Env
	PlacesService *places.Service
	fieldMask     string
	apiKey        string
	serverAddr    string
}

func NewMaps(env *env.Env) MapsIf {
	ctx := context.Background()
	placesService, err := places.NewService(ctx, option.WithAPIKey(env.GoogleMapsApiKey))
	if err != nil {
		panic(err)
	}

	return &Maps{
		PlacesService: placesService,
		env:           env,
		fieldMask:     "routes.legs.distanceMeters",
		apiKey:        env.GoogleMapsApiKey,
		serverAddr:    "routes.googleapis.com:443",
	}
}

func (m *Maps) GetAutocomplete(req dto.LocationRequest) ([]dto.LocationResponse, error) {
	if req.Lat == 0 {
		req.Lat = m.env.DefaultUserLatitude
	}

	if req.Long == 0 {
		req.Long = m.env.DefaultUserLongitude
	}

	if req.Radius == 0 {
		req.Radius = m.env.DefaultUserRadius
	}

	var suggestions []dto.LocationResponse

	request := &places.GoogleMapsPlacesV1AutocompletePlacesRequest{
		Input:        req.Query,
		RegionCode:   "ID",
		LocationBias: &places.GoogleMapsPlacesV1AutocompletePlacesRequestLocationBias{},
	}

	request.LocationBias.Circle = &places.GoogleMapsPlacesV1Circle{
		Center: &places.GoogleTypeLatLng{
			Latitude:  req.Lat,
			Longitude: req.Long,
		},
		Radius: req.Radius,
	}

	response, err := m.PlacesService.Places.Autocomplete(request).Do()
	if err != nil {
		log.Println("Error fetching autocomplete:", err)
		return suggestions, err
	}

	for _, suggestion := range response.Suggestions {
		response, _ := m.PlacesService.Places.Get("places/" + suggestion.PlacePrediction.PlaceId).Fields("*").Do()
		suggestions = append(suggestions, dto.LocationResponse{
			Name:    response.DisplayName.Text,
			PlaceID: suggestion.PlacePrediction.PlaceId,
		})
	}

	return suggestions, nil
}

func (m *Maps) GetPlaceDetails(placeId string) (dto.PlaceDetails, error) {
	response, err := m.PlacesService.Places.Get("places/" + placeId).Fields("*").Do()
	if err != nil {
		log.Println("Error fetching place details:", err)
		return dto.PlaceDetails{}, err
	}

	latitude := response.Location.Latitude
	longitude := response.Location.Longitude

	return dto.PlaceDetails{
		Name:    response.DisplayName.Text,
		PlaceId: placeId,
		Lat:     latitude,
		Long:    longitude,
	}, nil
}

func (m *Maps) GetDistance(from dto.Location, to dto.Location) (*int32, error) {
	config := tls.Config{}
	conn, err := grpc.NewClient(m.serverAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(&config)))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := routespb.NewRoutesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-Api-Key", m.apiKey)
	ctx = metadata.AppendToOutgoingContext(ctx, "X-Goog-Fieldmask", m.fieldMask)
	defer cancel()

	// create the origin using a latitude and longitude
	origin := &routespb.Waypoint{
		LocationType: &routespb.Waypoint_Location{
			Location: &routespb.Location{
				LatLng: &latlng.LatLng{
					Latitude:  from.Lat,
					Longitude: from.Long,
				},
			},
		},
	}

	// create the destination using a latitude and longitude
	destination := &routespb.Waypoint{
		LocationType: &routespb.Waypoint_Location{
			Location: &routespb.Location{
				LatLng: &latlng.LatLng{
					Latitude:  to.Lat,
					Longitude: to.Long,
				},
			},
		},
	}

	req := &routespb.ComputeRoutesRequest{
		Origin:      origin,
		Destination: destination,
		RouteModifiers: &routespb.RouteModifiers{
			AvoidTolls:    true,
			AvoidHighways: false,
			AvoidFerries:  true,
		},
		TravelMode: routespb.RouteTravelMode_TWO_WHEELER,
	}

	resp, err := client.ComputeRoutes(ctx, req)

	if err != nil {
		return nil, err
	}

	distance := resp.GetRoutes()[0].GetLegs()[0].GetDistanceMeters()

	return &distance, nil
}

func (m *Maps) GenerateGoogleMapsURL(placeId string) string {
	return "https://www.google.com/maps/place?q=place_id:" + placeId
}
