# Minimal CLI Example

A simple command-line interface demonstrating the go-agent-runner framework with tool support and both interactive and non-interactive modes.

## Features

- **Interactive Mode**: REPL-style chat session with command support
- **Non-Interactive Mode**: One-shot queries via `--prompt` flag
- **Tool Support**: Weather, locations, time, and batch weather tools
- **Smart Routing**: LLM-based directive selection (RouterMultiplexer)
- **Session Persistence**: Save and load chat sessions

## Build

```bash
go build ./...
```

## Run

### Interactive Mode (REPL)

```bash
go run cmd/main.go
```

Once in the REPL, you can:
- Type messages to chat with the agent
- Use `/help` for available commands
- Use `/save` to save the current session
- Use `/load` to load a previous session
- Use `/exit` or `/quit` to exit

### Non-Interactive Mode

```bash
# Simple query
go run cmd/main.go --prompt "what time is it?"

# Complex query with tools and formatting
go run cmd/main.go --prompt "get the current time. then the weather for munich and paris. put the time and the weather in a CSV"
```

## Available Tools

| Tool | Description |
|------|-------------|
| `get_current_time` | Get the current time |
| `get_locations` | Get coordinates for major world cities |
| `get_weather` | Get weather for a single location |
| `get_weather_batch` | Get weather for multiple locations at once |

## Directives

The CLI uses a RouterMultiplexer that intelligently selects which directives (and their associated tools) to activate based on the user's request:

- **weather**: Activates `get_weather` tool
- **weather-batch**: Activates `get_weather_batch` tool for multi-location queries
- **locations**: Activates `get_locations` tool
- **time**: Activates `get_current_time` tool
- **table-format**: Formats output as Markdown tables
- **csv-format**: Formats output as CSV

## Configuration

Flags:
- `-url`: OpenAI-compatible API base URL (default: `http://10.0.1.114:1234`)
- `-model`: Model to use (default: `mistralai/ministral-3-3b`)
- `-prompt`: Non-interactive mode prompt (empty = interactive mode)

## Example Session

```
 Replikant Agent - Minimal CLI Example

Type your message or /help for commands.

> whats the weather in tokyo? Format the result as a table.
╔Routing══════════
╟─[options]: [locations weather table-format]
║  - User requests weather information specifically for Tokyo, which
║ falls under single-location weather queries. - To retrieve weather
║ data for Tokyo, we need `get_locations` to first obtain the
║ coordinates (latitude/longitude) of Tokyo. - Once coordinates are
║ retrieved, `get_weather` will provide current weather details. - The
║ user also explicitly requested a table format via `<table-format>`. -
║ **Relevant directives**: - `locations`: Required to fetch Tokyo's
║ coordinates before calling `get_weather`. - `weather`: Core directive
║ for retrieving weather data (requires coordinates from `locations`). -
║ `table-format`: Ensures output is formatted as a Markdown table.
╚═════════════════

╭Directives─────────────────────────────────────────────
│╭[locations]:
││├─Prompt: 'USE get_locations when you need to find ...'
││╰─Tools: [get_locations]
│├[weather]:
││├─Prompt: ' WEATHER QUERY WORKFLOWS:  For single-lo...'
││╰─Tools: [get_weather get_weather_batch]
│├[table-format]:
││├─Prompt: 'Format the output as a Markdown table. W...'
││╰─Tools: ✗
╰╯

┌──────────────────────────────────────────────────────────
├ Tool call: get_locations({})
└─────────────────────────────────────────────────────────
  ✓ Available locations: 1. Tokyo, Japan - Lat: 35.6762, Lon: 139.6503 2. Delhi, India - Lat: 28.7041, Lon: 77.1025 3. Shanghai, China - Lat: 31.2304, Lon: 121.4737 4. São Paulo, Brazil - Lat: -23.5505, ...
┌──────────────────────────────────────────────────────────
├ Tool call: get_weather({"latitude":35.6762,"longitude":139.6503})
└─────────────────────────────────────────────────────────
  ✓ Current weather: 15.0°C, Light drizzle, wind: 6.3 km/h

Agent: 
| CITY        | TEMPERATURE | CONDITIONS          | WIND SPEED |
|-------------|-------------|---------------------|------------|
| TOKYO       | 15.0°C      | LIGHT DRIZZLE       | 6.3 KM/H   |
```
