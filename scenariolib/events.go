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
	switch e.Type {

	case "Click":
		clickEvent := &ClickEvent{}
		if err := json.Unmarshal(e.Arguments, clickEvent); err != nil {
			return nil, err
		}
		if valid, message := clickEvent.IsValid(); !valid {
			return nil, errors.New(message)
		}
		return clickEvent, nil

	case "Custom":
		customEvent := &CustomEvent{}
		if err := json.Unmarshal(e.Arguments, customEvent); err != nil {
			return nil, err
		}
		if valid, message := customEvent.IsValid(); !valid {
			return nil, errors.New(message)
		}
		return customEvent, nil

	case "FacetChange":
		facetEvent := &FacetEvent{}
		if err := json.Unmarshal(e.Arguments, facetEvent); err != nil {
			return nil, err
		}
		if valid, message := facetEvent.IsValid(); !valid {
			return nil, errors.New(message)
		}
		return facetEvent, nil

	case "Search":
		event, err := newSearchEvent(e, c)
		if err != nil {
			return nil, err
		}
		return event, nil

	case "FakeSearch":
		event, err := newFakeSearchEvent(e, c)
		if err != nil {
			return nil, err
		}
		return event, nil

	case "SearchAndClick":
		event, err := newSearchAndClickEvent(e)
		if err != nil {
			return nil, err
		}
		return event, nil

	case "TabChange":
		event, err := newTabChangeEvent(e)
		if err != nil {
			return nil, err
		}
		return event, nil
	case "View":
		event, err := newViewEvent(e, c)
		if err != nil {
			return nil, err
		}
		return event, nil

	case "SetOrigin":
		event, err := newSetOriginEvent(e)
		if err != nil {
			return nil, err
		}
		return event, nil

	case "SetReferrer":
		event, err := newSetReferrerEvent(e)
		if err != nil {
			return nil, err
		}
		return event, nil
	}
	return nil, errors.New("Event type not supported")
}

// Event Generic interface for abstract type Event. All specific event types must
// define the Execute function
type Event interface {
	Execute(v *Visit) error
	IsValid() (bool, string)
}
