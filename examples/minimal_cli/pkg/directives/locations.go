package directives

import (
	"github.com/sjsone/go-replikant-agent/examples/minimal_cli/pkg/tools/locations"
	"github.com/sjsone/go-replikant-agent/lib/directive"
	"github.com/sjsone/go-replikant-agent/lib/prompt"
	"github.com/sjsone/go-replikant-agent/lib/tool"
)

// NewLocationsDirective creates a locations lookup directive
func NewLocationsDirective() directive.Directive {
	return directive.NewStaticDirective(
		"locations",
		&prompt.Prompt{Raw: "USE get_locations when you need to find coordinates (latitude/longitude) for a city. Returns a list of major world cities with their coordinates."},
		[]tool.ToolCallable{locations.NewLocationsTool()},
	)
}
