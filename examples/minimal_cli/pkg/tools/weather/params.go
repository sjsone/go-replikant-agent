package weather

// WeatherParams defines parameters for get_weather tool.
type WeatherParams struct {
	Latitude  float64 `json:"latitude" jsonschema:"Latitude of the location (e.g., 51.5074 for London)"`
	Longitude float64 `json:"longitude" jsonschema:"Longitude of the location (e.g., -0.1278 for London)"`
}

// WeatherBatchParams defines parameters for get_weather_batch tool.
type WeatherBatchParams struct {
	Locations []WeatherParams `json:"locations" jsonschema:"Array of locations to get weather for"`
}
