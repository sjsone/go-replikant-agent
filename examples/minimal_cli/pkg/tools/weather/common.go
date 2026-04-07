package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func weatherCodeToDescription(code int) string {
	codes := map[int]string{
		0:  "Clear sky",
		1:  "Mainly clear",
		2:  "Partly cloudy",
		3:  "Overcast",
		45: "Foggy",
		48: "Depositing rime fog",
		51: "Light drizzle",
		53: "Moderate drizzle",
		55: "Dense drizzle",
		61: "Slight rain",
		63: "Moderate rain",
		65: "Heavy rain",
		71: "Slight snow",
		73: "Moderate snow",
		75: "Heavy snow",
		80: "Slight rain showers",
		81: "Moderate rain showers",
		82: "Violent rain showers",
		95: "Thunderstorm",
	}
	if desc, ok := codes[code]; ok {
		return desc
	}
	return "Unknown weather code"
}

// fetchWeatherForLocation fetches weather data for a single location.
func fetchWeatherForLocation(ctx context.Context, loc WeatherParams) (string, error) {
	apiURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,weather_code,wind_speed_10m",
		loc.Latitude, loc.Longitude,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("weather API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		Current struct {
			Temperature float64 `json:"temperature_2m"`
			WeatherCode int     `json:"weather_code"`
			WindSpeed   float64 `json:"wind_speed_10m"`
		} `json:"current"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse weather data: %w", err)
	}

	return fmt.Sprintf(
		"%.1f°C, %s, wind: %.1f km/h",
		result.Current.Temperature,
		weatherCodeToDescription(result.Current.WeatherCode),
		result.Current.WindSpeed,
	), nil
}

// locationResult holds the result for a single location.
type locationResult struct {
	Name    string
	Weather string
	Err     error
}
