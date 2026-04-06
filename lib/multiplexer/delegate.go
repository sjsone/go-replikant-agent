package multiplexer

import "github.com/sjsone/go-replikant-agent/lib/directive"

type Delegate interface {
	// MultiplexerOnDirectivesSelected is called after the multiplexer selects active directives.
	// Provides both the full set of available directives and the selected subset.
	//
	// allDirectives: All directives registered with the multiplexer
	// activeDirectives: The directives selected for this request
	MultiplexerOnDirectivesSelected(allDirectives []directive.Directive, activeDirectives []directive.Directive)
}
