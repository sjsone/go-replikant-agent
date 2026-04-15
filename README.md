# Replikant Agent

> [!IMPORTANT] Evaluation Process
> This is a WIP project evaluating an agentic framework focused on the implementation of system directives and modular prompt architecture.

```
╭────────────╮                      
│ User Query │                      
╰───────┬────╯        ┌─────────────┐
        │          ┌──┤ Directive A │
┌───────▼─────┐    │  └─────────────┘
│ Multiplexer ◄────┤  ┌┈┈┈┈┈┈┈┈┈┈┈┈┈┐              
└───────┬─────┘    └┈┈┤ Directive B │            
        │             └┈┈┈┈┈┈┈┈┈┈┈┈┈┘            
╔═══════▼═══════╗                 
║ Loop          ◄┐                
╚═══════╤═══════╝│                 
        │        │                  
┌───────▼────┐   │                  
│    LLM     │   │                  
└┬──────┬────┘   │                  
 │      │        │                  
 │ ┌────▼──────┐ │                  
 │ │ Tool Call ├─┘                  
 │ └───────────┘                    
 │                                
╭▼─────────╮                      
│ Response │
╰──────────╯
```

Replikant Agent is a modular Go framework for building LLM-powered agents with pluggable connectors, directives (capability modules with prompts and tools), routing, and session management.

## Quick start

See [`examples/minimal_cli/`](examples/minimal_cli/) for a working CLI that demonstrates interactive and one-shot modes, tool execution, and directive routing.
