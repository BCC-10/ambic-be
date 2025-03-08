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

type Place struct {
	lat  float64
	long float64
}

type MapsIf interface {
	GetAutocomplete(req dto.LocationRequest) ([]dto.LocationResponse, error)
	GetPlaceDetails(placeId string) (PlaceDetails, error)
	GetDistance(from PlaceDetails, to PlaceDetails) any
	GenerateGoogleMapsURL(placeId string) string
}

type Maps struct {
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
		fieldMask:     "routes.legs.distanceMeters",
		apiKey:        env.GoogleMapsApiKey,
		serverAddr:    "routes.googleapis.com:443",
	}
}

func (m *Maps) GetAutocomplete(req dto.LocationRequest) ([]dto.LocationResponse, error) {
	if req.Radius == 0 {
		req.Radius = 20000
	}

	var suggestions []dto.LocationResponse

	request := &places.GoogleMapsPlacesV1AutocompletePlacesRequest{
		Input:        req.Query,
		RegionCode:   "ID",
		LocationBias: &places.GoogleMapsPlacesV1AutocompletePlacesRequestLocationBias{},
	}
	if req.Lat != 0 || req.Long != 0 {
		req.Lat = -6.175110
		req.Long = 106.865036

		request.LocationBias.Circle = &places.GoogleMapsPlacesV1Circle{
			Center: &places.GoogleTypeLatLng{
				Latitude:  req.Lat,
				Longitude: req.Long,
			},
			Radius: req.Radius,
		}
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

type PlaceDetails struct {
	Name    string
	PlaceId string
	Lat     float64
	Long    float64
}

func (m *Maps) GetPlaceDetails(placeId string) (PlaceDetails, error) {
	response, err := m.PlacesService.Places.Get("places/" + placeId).Fields("*").Do()
	if err != nil {
		log.Println("Error fetching place details:", err)
		return PlaceDetails{}, err
	}

	latitude := response.Location.Latitude
	longitude := response.Location.Longitude

	return PlaceDetails{
		Name:    response.DisplayName.Text,
		PlaceId: placeId,
		Lat:     latitude,
		Long:    longitude,
	}, nil
}

func (m *Maps) GetDistance(from PlaceDetails, to PlaceDetails) any {
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
			AvoidFerries:  false,
		},
	}

	resp, err := client.ComputeRoutes(ctx, req)

	if err != nil {
		// "rpc error: code = InvalidArgument desc = Request contains an invalid
		// argument" may indicate that your project lacks access to Routes
		log.Fatal(err)
	}

	return resp
}

func (m *Maps) GenerateGoogleMapsURL(placeId string) string {
	return "https://www.google.com/maps/place?q=" + placeId
}
