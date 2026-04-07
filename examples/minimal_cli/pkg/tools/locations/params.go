package locations

// LocationsParams defines parameters for get_locations tool (empty).
type LocationsParams struct{}

// LocationInput defines a single location in the batch request.
type LocationInput struct {
	Name      string  `json:"name" jsonschema:"Name of the location/city for display"`
	Latitude  float64 `json:"latitude" jsonschema:"Latitude of the location"`
	Longitude float64 `json:"longitude" jsonschema:"Longitude of the location"`
}
