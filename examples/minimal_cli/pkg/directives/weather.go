package directives

import (
	"github.com/sjsone/go-replikant-agent/examples/minimal_cli/pkg/tools/weather"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// NewWeatherDirective creates a unified weather query directive with all weather-related tools
func NewWeatherDirective() directive.Directive {
	return directive.NewStaticDirective(
		"weather",
		&prompt.Prompt{Raw: `
WEATHER QUERY WORKFLOWS:

For single-location weather queries:
1. Call get_locations to see available cities and their coordinates
2. Call get_weather with the desired city's latitude and longitude

For multi-location weather queries:
- Use get_weather_batch with an array of locations (each with name, latitude, longitude)
- This is more efficient than multiple individual get_weather calls

Example multi-location format:
[
  {"name": "Tokyo", "latitude": 35.6762, "longitude": 139.6503},
  {"name": "New York", "latitude": 40.7128, "longitude": -74.0060}
]
`},
		[]tool.ToolCallable{
			weather.NewWeatherTool(),
			weather.NewWeatherBatchTool(),
		},
	)
}
