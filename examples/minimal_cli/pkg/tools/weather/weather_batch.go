package weather

import (
	"context"
	"fmt"
	"sync"

	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// WeatherBatchTool implements the ToolCallable interface for batch weather queries.
type WeatherBatchTool struct {
	tool.Tool
}

// NewWeatherBatchTool creates a new batch weather tool.
func NewWeatherBatchTool() tool.ToolCallable {
	w := &WeatherBatchTool{}

	params, err := tool.SchemaFromStruct[WeatherBatchParams]()
	if err != nil {
		panic(err)
	}

	w.Tool = tool.Tool{
		Name:        "get_weather_batch",
		Description: "Get current weather information for multiple locations at once. Takes an array of locations with name, latitude, and longitude. Returns weather data for all locations efficiently.",
		Parameters:  params,
	}
	return w
}

// GetTool returns the tool metadata.
func (w *WeatherBatchTool) GetTool() *tool.Tool {
	return &w.Tool
}

// Execute implements ToolCallable with concurrent HTTP requests.
func (w *WeatherBatchTool) Execute(ctx context.Context, args map[string]any) (string, error) {
	locationsRaw, ok := args["locations"].([]any)
	if !ok {
		return "", fmt.Errorf("invalid locations parameter: expected array")
	}

	if len(locationsRaw) == 0 {
		return "", fmt.Errorf("locations array cannot be empty")
	}

	// Parse location inputs
	locationInputs := make([]WeatherParams, 0, len(locationsRaw))
	for i, locRaw := range locationsRaw {
		locMap, ok := locRaw.(map[string]any)
		if !ok {
			return "", fmt.Errorf("invalid location at index %d: expected object", i)
		}

		lat, ok_lat := locMap["latitude"].(float64)
		lon, ok_lon := locMap["longitude"].(float64)
		if !ok_lat || !ok_lon {
			return "", fmt.Errorf("invalid location at index %d: missing required fields", i)
		}

		locationInputs = append(locationInputs, WeatherParams{
			Latitude:  lat,
			Longitude: lon,
		})
	}

	// Fetch weather concurrently
	results := make([]locationResult, len(locationInputs))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for i, loc := range locationInputs {
		wg.Add(1)
		go func(idx int, location WeatherParams) {
			defer wg.Done()

			// Check for context cancellation before starting request
			select {
			case <-ctx.Done():
				mu.Lock()
				if firstErr == nil {
					firstErr = ctx.Err()
				}
				mu.Unlock()
				return
			default:
			}

			weather, err := fetchWeatherForLocation(ctx, location)
			mu.Lock()
			results[idx] = locationResult{
				Weather: weather,
				Err:     err,
			}
			if err != nil && firstErr == nil {
				firstErr = err
			}
			mu.Unlock()
		}(i, loc)
	}

	wg.Wait()

	// If all requests failed due to context cancellation, return the error
	if firstErr == context.Canceled || firstErr == context.DeadlineExceeded {
		return "", firstErr
	}

	// Build output string
	var output string
	for i, result := range results {
		// Check for context cancellation during output building
		select {
		case <-ctx.Done():
			return output, ctx.Err()
		default:
		}

		if result.Err != nil {
			output += fmt.Sprintf("%d. %s: Error - %v\n", i+1, result.Name, result.Err)
		} else {
			output += fmt.Sprintf("%d. %s: %s\n", i+1, result.Name, result.Weather)
		}
	}

	return output, nil
}
