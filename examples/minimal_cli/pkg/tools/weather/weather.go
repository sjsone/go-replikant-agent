package weather

import (
	"context"
	"fmt"

	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// WeatherTool implements the ToolCallable interface for weather queries.
type WeatherTool struct {
	tool.Tool
}

// NewWeatherTool creates a new weather tool.
func NewWeatherTool() tool.ToolCallable {
	w := &WeatherTool{}

	params, err := tool.SchemaFromStruct[WeatherParams]()
	if err != nil {
		panic(err)
	}

	w.Tool = tool.Tool{
		Name:        "get_weather",
		Description: "Get current weather information for a location",
		Parameters:  params,
	}
	return w
}

// GetTool returns the tool metadata.
func (w *WeatherTool) GetTool() *tool.Tool {
	return &w.Tool
}

func (w *WeatherTool) GetName() string {
	return w.Tool.Name
}

// Execute implements ToolCallable.
func (w *WeatherTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	lat, ok1 := args["latitude"].(float64)
	lon, ok2 := args["longitude"].(float64)
	if !ok1 || !ok2 {
		return "", fmt.Errorf("invalid latitude or longitude")
	}

	weather, err := fetchWeatherForLocation(ctx, WeatherParams{Latitude: lat, Longitude: lon})
	if err != nil {
		return "", err
	}

	return "Current weather: " + weather, nil
}
