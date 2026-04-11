VERSION 0.1

# Replikant Agent Framework

## Overview

# Replikant Agent is a modular Go framework for building LLM-powered agents. It provides pluggable connectors, directives (capability modules with prompts and tools), routing, and session management.

## Architecture

# The framework follows a pipeline:
# User Input → Multiplexer (selects directives) → Session (builds prompts, calls LLM, executes tools) → Loop Controller (prevents infinite loops)

## Type System Reference

# BHaus | Description
# -------|-------------
# String | String type
# Integer | Integer type
# Boolean | Boolean type
# Unknown | Any/unknown type
# ?Type | Nullable/optional type
# Array[Type] | Array of Type

## Core Protocols

### Connector
# Interface for LLM API communication with streaming support.

PROTOCOL Connector:
    PUBLIC Send(context, messages: Array[Message], directives: Array[Directive], onChunk: ChunkHandler): (error, ?ContextPart)

### RoutingConnector
# Interface for routing-specific requests that return raw JSON.

PROTOCOL RoutingConnector:
    PUBLIC SendForRouting(context, messages: Array[ChatMessage]): (raw JSON, error)

### Directive
# Bundles a prompt, tool definitions, and tool callables into a composable capability module.

PROTOCOL Directive:
    PUBLIC GetName(): String
    PUBLIC GetPrompt(): Prompt
    PUBLIC GetTools(): Array[Tool]
    PUBLIC GetToolCallables(): Array[ToolCallable]

### Multiplexer
# Selects which directives are active for a given context.

PROTOCOL Multiplexer:
    PUBLIC GetActiveDirectivesForContext(context: AgentContext): Array[Directive]
    PUBLIC GetAllDirectives(): Array[Directive]

### Router
# LLM-based directive selection. Takes a user query and available options, returns which directives to activate.

PROTOCOL Router:
    PUBLIC Route(context, userQuery: String, allAvailableOptions: Array[RoutingOption]): RoutingResult

### Session
# Orchestrates the core agent loop — builds messages, sends to connector, handles tool calls, manages context parts.

PROTOCOL Session:
    PUBLIC ProcessContextPart(context, part: ContextPart): error
    PUBLIC CurrentContext(): AgentContext
    PUBLIC Cancel()

### LoopController
# Determines when the agent loop should continue (e.g., when there are pending tool results to process).

PROTOCOL LoopController:
    PUBLIC LoopAgain(context: AgentContext): Boolean

### PromptBuilder
# Constructs prompts from directives.

PROTOCOL PromptBuilder:
    PUBLIC Build(directives: Array[Directive]): Prompt

### ToolCallable
# Interface tools must implement for execution.

PROTOCOL ToolCallable:
    PUBLIC Execute(context, args: Unknown): (String, error)
    PUBLIC GetName(): String

---

## Data Structures

### Message
# Represents a message with source and text content.

STRUCT Message:
    PUBLIC Source: String
    PUBLIC Text: String

### ChatMessage
# Generic chat message format used for routing requests.

STRUCT ChatMessage:
    PUBLIC Role: String
    PUBLIC Content: String
    PUBLIC ToolCallID: String
    PUBLIC ToolCalls: Array[FunctionCall]

### Prompt
# Simple wrapper around raw string prompt content.

STRUCT Prompt:
    PUBLIC Raw: String

### Tool
# Metadata for LLM-defined tools including name, description, and JSON schema parameters.

STRUCT Tool:
    PUBLIC Name: String
    PUBLIC Description: String
    PUBLIC Parameters: Unknown

### FunctionCall
# Parsed tool call from LLM response.

STRUCT FunctionCall:
    PUBLIC Name: String
    PUBLIC Arguments: Unknown

### FunctionResult
# Result of tool execution.

STRUCT FunctionResult:
    PUBLIC Name: String
    PUBLIC Result: String
    PUBLIC Error: Boolean

### ContextPart
# Single message/tool call/tool result in conversation history.

STRUCT ContextPart:
    PUBLIC Raw: String
    PUBLIC Source: String
    PUBLIC ToolUse: Boolean
    PUBLIC Stop: Boolean
    PUBLIC Cancelled: Boolean
    PUBLIC ToolCalls: Array[FunctionCall]
    PUBLIC ToolResults: Array[FunctionResult]
    PUBLIC ConnectedToolCallContextPart: ?ContextPart
    PUBLIC ConnectedToolCallSourceContextPart: ?ContextPart

### AgentContext
# Holds conversation history as a list of ContextParts.

STRUCT AgentContext:
    PUBLIC Parts: Array[ContextPart]

### RoutingOption
# Represents a directive option available for routing selection.

STRUCT RoutingOption:
    PUBLIC ID: String
    PUBLIC Name: String
    PUBLIC Description: String

### RoutingResult
# Contains the selected options and routing decision.

STRUCT RoutingResult:
    PUBLIC SelectedIDs: Array[String]
    PUBLIC Reasoning: String
    PUBLIC Confidence: Integer

### RoutingDecision
# The routing decision including selected IDs and reasoning.

STRUCT RoutingDecision:
    PUBLIC SelectedIDs: Array[String]
    PUBLIC Reasoning: String
    PUBLIC Confidence: Integer

### LoopDecision
# Decision about whether the agent loop should continue.

STRUCT LoopDecision:
    PUBLIC Continue: Boolean
    PUBLIC Reason: String

---

## Concrete Implementations

### StaticDirective
# Standard Directive implementation that stores prompt and tool callables directly.

CLASS StaticDirective IMPLEMENTS Directive:
    PUBLIC GetName(): String
    PUBLIC GetPrompt(): Prompt
    PUBLIC GetTools(): Array[Tool]
    PUBLIC GetToolCallables(): Array[ToolCallable]

### SimpleMultiplexer
# Activates all directives regardless of context.

CLASS SimpleMultiplexer IMPLEMENTS Multiplexer:
    PUBLIC GetActiveDirectivesForContext(context: AgentContext): Array[Directive]
    PUBLIC GetAllDirectives(): Array[Directive]

### RouterMultiplexer
# Uses LLM-based routing to select which directives are active.

CLASS RouterMultiplexer IMPLEMENTS Multiplexer:
    PUBLIC GetActiveDirectivesForContext(context: AgentContext): Array[Directive]
    PUBLIC GetAllDirectives(): Array[Directive]
    PUBLIC GetLastRoutingDecision(): RoutingDecision

### SimpleRouter
# Default LLM-based router implementation.

CLASS SimpleRouter IMPLEMENTS Router:
    PUBLIC Route(context, userQuery: String, allAvailableOptions: Array[RoutingOption]): RoutingResult
    PUBLIC SetExampleMessages(messages: Array[ChatMessage])
    PUBLIC SetDelegate(delegate: RouterDelegate)

### SimpleLoopController
# Default implementation that decides loop continuation based on context analysis.

CLASS SimpleLoopController IMPLEMENTS LoopController:
    PUBLIC LoopAgain(context: AgentContext): Boolean
    PUBLIC SetDelegate(delegate: LoopDelegate)

### AgenticSession
# Orchestrates the core agent loop.

CLASS AgenticSession:
    PUBLIC ProcessContextPart(context, part: ContextPart): error
    PUBLIC CurrentContext(): AgentContext
    PUBLIC Cancel()
    PUBLIC SetDelegate(delegate: SessionDelegate)

---

## Delegate Interfaces
# The framework uses the **delegate pattern** throughout for extensibility. All delegates are nil-safe — implementations can choose which methods to override.

### SessionDelegate
# Observer for session events.

PROTOCOL SessionDelegate:
    PUBLIC SessionOnPartAdded(part: ContextPart)
    PUBLIC SessionOnToolCallsReceived(calls: Array[FunctionCall])
    PUBLIC SessionOnToolExecuted(call: FunctionCall, result: FunctionResult)
    PUBLIC SessionOnStreamingChunk(chunk: String)
    PUBLIC SessionOnRequestSent(messages: Array[Message], directives: Array[Directive])
    PUBLIC SessionOnLoopIteration(iteration: Integer)
    PUBLIC SessionOnLoopEnd()

### RouterDelegate
# Observer for routing decisions.

PROTOCOL RouterDelegate:
    PUBLIC RouterOnRoutingDecision(decision: RoutingDecision, allOptions: Array[RoutingOption], activeOptions: Array[RoutingOption])

### LoopDelegate
# Observer for loop continuation decisions.

PROTOCOL LoopDelegate:
    PUBLIC LoopOnDecision(decision: LoopDecision)

## Notes
# - All interfaces follow the delegate pattern for extensibility
# - The framework is designed to be lightweight with minimal external dependencies
# - The only external dependency is github.com/google/jsonschema-go for JSON schema validation
# - Streaming is supported through the ChunkHandler callback
# - The loop controller prevents infinite loops by analyzing context and tool results
# - Routing can be done via simple multiplexing (all active) or LLM-based selection
