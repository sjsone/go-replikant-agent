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
> get the weather for Tokyo
[Directives 2/6 active]
│ locations, weather
[Routing indices=[2 0]]
│ User wants weather for Tokyo. Need to first get Tokyo's coordinates...

┌──────────────────────────────────────────────────────────
│ Tool call: get_locations({})
└─────────────────────────────────────────────────────────
  ✓ Available locations: 1. Tokyo, Japan - Lat: 35.6762...

┌──────────────────────────────────────────────────────────
│ Tool call: get_weather({"latitude":35.6762,"longitude":139.6503})
└─────────────────────────────────────────────────────────
  ✓ Tokyo: 15.3°C, Cloudy, wind: 12.1 km/h

Agent: The current weather in Tokyo is 15.3°C with cloudy skies and wind at 12.1 km/h.
```
