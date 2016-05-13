// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import "errors"

// ============== SEARCH EVENT ======================
// ==================================================

// CustomEvent a struct representing a search, is defined by a query to execute
type CustomEvent struct {
	eventType  string
	eventValue string
}

func newCustomEvent(e *JSONEvent, c *Config) (*CustomEvent, error) {
	eventType, ok1 := e.Arguments["eventType"].(string)
	eventValue, ok2 := e.Arguments["eventValue"].(string)
	if !ok1 || !ok2 {
		return nil, errors.New("Invalid parse of arguments on Custom Event")
	}

	return &CustomEvent{
		eventType:  eventType,
		eventValue: eventValue,
	}, nil
}

// Execute Execute the search event, runs the query and sends a search event to
// the analytics.
func (ce *CustomEvent) Execute(v *Visit) error {
	Info.Printf("CustomEvent type: %s ||| value: %s", ce.eventType, ce.eventValue)

	err := v.sendCustomEvent(ce.eventType, ce.eventValue)
	return err
}
