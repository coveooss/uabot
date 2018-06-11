// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import (
	"encoding/json"
	"errors"
)

// ParseEvent A factory to create the correct event type coming from the JSON parse
// of the scenario definition.
func ParseEvent(e *JSONEvent, c *Config) (Event, error) {

	var event Event
	switch e.Type {

	case "Click":
		event = &ClickEvent{}

	case "Custom":
		event = &CustomEvent{}

	case "FacetChange":
		event = &FacetEvent{}

	case "FakeSearch":
		event = &FakeSearchEvent{}

	case "View":
		event = &ViewEvent{}

	case "Search":
		event = &SearchEvent{}

	case "SearchAndClick":
		event = &SearchAndClickEvent{}

	case "SetOrigin":
		event = &SetOriginEvent{}

	case "SetReferrer":
		event = &SetReferrerEvent{}

	case "TabChange":
		event = &TabChangeEvent{}

	default:
		return nil, errors.New("Event type not supported")

	}

	if err := json.Unmarshal(e.Arguments, event); err != nil {
		return nil, err
	}
	if valid, message := event.IsValid(); !valid {
		return nil, errors.New(message)
	}
	return event, nil
}

// Event Generic interface for abstract type Event. All specific event types must
// define the Execute function
type Event interface {
	Execute(v *Visit) error
	IsValid() (bool, string)
}
