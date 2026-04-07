package locations

import (
	"context"
	"fmt"
	"strings"

	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// Location represents a city with its coordinates.
type Location struct {
	Name      string
	Country   string
	Latitude  float64
	Longitude float64
}

// LocationsTool implements the ToolCallable interface for location lookup.
type LocationsTool struct {
	tool.Tool
}

// Top 20 cities by population with coordinates
var cities = []Location{
	{"Tokyo", "Japan", 35.6762, 139.6503},
	{"Delhi", "India", 28.7041, 77.1025},
	{"Shanghai", "China", 31.2304, 121.4737},
	{"São Paulo", "Brazil", -23.5505, -46.6333},
	{"Mexico City", "Mexico", 19.4326, -99.1332},
	{"Cairo", "Egypt", 30.0444, 31.2357},
	{"Mumbai", "India", 19.0760, 72.8777},
	{"Beijing", "China", 39.9042, 116.4074},
	{"Dhaka", "Bangladesh", 23.8103, 90.4125},
	{"Osaka", "Japan", 34.6937, 135.5023},
	{"New York", "USA", 40.7128, -74.0060},
	{"Karachi", "Pakistan", 24.8607, 67.0011},
	{"Buenos Aires", "Argentina", -34.6037, -58.3816},
	{"Chongqing", "China", 29.4316, 106.9123},
	{"Istanbul", "Turkey", 41.0082, 28.9784},
	{"Kolkata", "India", 22.5726, 88.3639},
	{"Manila", "Philippines", 14.5995, 120.9842},
	{"Lagos", "Nigeria", 6.5244, 3.3792},
	{"Rio de Janeiro", "Brazil", -22.9068, -43.1729},
	{"London", "UK", 51.5074, -0.1278},
	{"Berlin", "Germany", 52.5200, 13.4050},
	{"Munich", "Germany", 48.1351, 11.5820},
	{"Paris", "France", 48.8566, 2.3522},
}

// NewLocationsTool creates a new locations tool.
func NewLocationsTool() tool.ToolCallable {
	l := &LocationsTool{}

	params, err := tool.SchemaFromStruct[LocationsParams]()
	if err != nil {
		panic(err)
	}

	l.Tool = tool.Tool{
		Name:        "get_locations",
		Description: "Get a list of major world cities with their latitude and longitude coordinates",
		Parameters:  params,
	}
	return l
}

// GetTool returns the tool metadata.
func (l *LocationsTool) GetTool() *tool.Tool {
	return &l.Tool
}

func (l *LocationsTool) GetName() string {
	return l.Tool.Name
}

// Execute implements ToolCallable.
func (l *LocationsTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	var sb strings.Builder

	sb.WriteString("Available locations:\n")
	for i, city := range cities {
		// Check for context cancellation during iteration
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}
		sb.WriteString(fmt.Sprintf("%d. %s, %s - Lat: %.4f, Lon: %.4f\n",
			i+1, city.Name, city.Country, city.Latitude, city.Longitude))
	}

	return sb.String(), nil
}
