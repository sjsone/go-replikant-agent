package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sjsone/go-replikant-agent/examples/minimal_cli/pkg/cli"
	"github.com/sjsone/go-replikant-agent/examples/minimal_cli/pkg/directives"
	"github.com/sjsone/go-replikant-agent/examples/minimal_cli/pkg/mode"
	"github.com/sjsone/go-replikant-agent/lib/agentic_context"
	"github.com/sjsone/go-replikant-agent/lib/connector/openai"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	simple_loop "github.com/sjsone/go-replikant-agent/lib/loop/simple"
	mcpclient "github.com/sjsone/go-replikant-agent/lib/mcp/client"
	mcpdirective "github.com/sjsone/go-replikant-agent/lib/mcp/directive"
	router_multiplexer "github.com/sjsone/go-replikant-agent/lib/multiplexer/router"
	prompt_string_builder "github.com/sjsone/go-replikant-agent/lib/prompt/builder/string"
	"github.com/sjsone/go-replikant-agent/lib/router/simple"
	lib_session "github.com/sjsone/go-replikant-agent/lib/session"
)

var (
	baseURLFlag  = flag.String("url", "", "OpenAI-compatible API base URL (or set URL_DEFAULT env)")
	modelFlag    = flag.String("model", "", "Model to use (or set MODEL_DEFAULT env)")
	promptArg    = flag.String("prompt", "", "Non-interactive mode: process prompt and exit")
	systemPrompt = flag.String("system-prompt", "", "System prompt for the agent (default: \"You are a helpful assistant\")")
	apiKeyFlag   = flag.String("api-key", "", "API key for authentication (or set API_KEY_DEFAULT env)")
	envFile      = flag.String("env", "", "Load environment from a .env file (default: .env if flag is present without value)")
)

func main() {
	// Handle bare --env (no value) before flag.Parse, since Go's flag package
	// requires an argument for non-boolean flags.
	args := os.Args[1:]
	for i, a := range args {
		if a == "--env" || a == "-env" {
			// If next arg looks like a flag or there is no next arg, treat as bare --env
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
				args[i] = "--env=.env"
			}
			break
		}
	}
	flag.CommandLine.Parse(args)

	// Load .env if --env flag is present
	if *envFile != "" {
		cli.LoadDotEnv(*envFile)
	}

	baseURL := cli.EnvOrFlag("URL_DEFAULT", *baseURLFlag)
	model := cli.EnvOrFlag("MODEL_DEFAULT", *modelFlag)
	apiKey := cli.EnvOrFlag("API_KEY_DEFAULT", *apiKeyFlag)

	if baseURL == "" || model == "" {
		fmt.Fprintln(os.Stderr, "Error: both --url and --model are required.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  Set via flags:  go run ./examples/minimal_cli/cmd/main.go --url <URL> --model <MODEL>")
		fmt.Fprintln(os.Stderr, "  Set via env:     URL_DEFAULT=<URL> MODEL_DEFAULT=<MODEL> go run ./examples/minimal_cli/cmd/main.go")
		fmt.Fprintln(os.Stderr, "  Set via .env:    --env (loads .env from working directory) or --env=/path/to/file")
		os.Exit(1)
	}

	sysPrompt := *systemPrompt
	if sysPrompt == "" {
		sysPrompt = "You are a helpful assistant"
	}

	systemPromptBuilder := prompt_string_builder.NewStringPromptBuilder(sysPrompt)

	loopController := simple_loop.NewSimpleLoopController()

	config := openai.DefaultOpenAIConfig(baseURL, model)
	config.APIKey = apiKey
	config.Timeout = 120 * time.Second
	con := openai.NewOpenAIConnector(config)

	agenticContext := agentic_context.NewAgentContext()

	allDirectives := []directive.Directive{
		directives.NewWeatherDirective(),
		directives.NewLocationsDirective(),
		directives.NewTimeDirective(),
		directives.NewTableFormattingDirective(),
		directives.NewCSVFormattingDirective(),
	}

	// Connect to AWS Knowledge MCP server.
	{
		config := mcpclient.ServerConfig{
			Name: "aws-knowledge",
			URL:  "https://knowledge-mcp.global.api.aws",
		}
		mcpDir, err := mcpdirective.NewMCPDirective(context.Background(), config)
		if err != nil {
			log.Printf("Warning: could not connect to AWS Knowledge MCP server: %v", err)
		} else {
			defer mcpDir.Close()
			allDirectives = append(allDirectives, mcpDir.Directive())
		}
	}

	router := simple.NewSimpleRouter("", con)

	multiplexer := router_multiplexer.NewRouterMultiplexer(allDirectives, router)
	// multiplexer := multiplexer.NewSimpleMultiplexer(allDirectives)

	agentic_session := lib_session.NewAgenticSession(multiplexer, *agenticContext, systemPromptBuilder, loopController, con)

	if *promptArg != "" {
		mode.RunSingleTurnMode(agentic_session, agenticContext, *promptArg, router)
		return
	}

	mode.RunInteractiveMode(agentic_session, agenticContext, router, &allDirectives)
}
