package router

// Delegate is an optional interface that delegates can implement
// to receive LLM routing decisions when using RouterMultiplexer.
type Delegate interface {
	RouterPreparedRouting(allOptions []*RoutingOption)

	// RouterOnRoutingDecision is called when RouterMultiplexer makes a routing decision.
	// Only called for RouterMultiplexer (not SimpleMultiplexer).
	//
	// decision: The routing decision containing selected IDs and reasoning
	// allOptions: All options available to the router
	// activeOptions: The options that were selected based on the decision
	RouterOnRoutingDecision(decision RoutingDecision, allOptions []*RoutingOption, activeOptions []*RoutingOption)
}
