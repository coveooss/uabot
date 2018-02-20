// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import "errors"

// ============== SEARCH EVENT ======================
// ==================================================

// CustomEvent a struct representing a search, is defined by a query to execute
type CustomEvent struct {
	actionCause string
	actionType  string
	customData  map[string]interface{}
}

func newCustomEvent(e *JSONEvent) (*CustomEvent, error) {
	var actionType, actionCause string
	var customData map[string]interface{}
	var validCast bool
	if actionType, validCast = e.Arguments["actionType"].(string); !validCast {
		return nil, errors.New("Parameter actionType is required and must be a string in a custom event")
	}
	if actionCause, validCast = e.Arguments["actionCause"].(string); !validCast {
		return nil, errors.New("Parameter actionCause is required and must be a string in a custom event")
	}
	if e.Arguments["customData"] != nil {
		if customData, validCast = e.Arguments["customData"].(map[string]interface{}); !validCast {
			return nil, errors.New("Parameter customData must be a json object (map[string]interface{}) in a custom event")
		}
	}

	return &CustomEvent{
		actionType:  actionType,
		actionCause: actionCause,
		customData:  customData,
	}, nil
}

// Execute Execute the search event, runs the query and sends a search event to
// the analytics.
func (ce *CustomEvent) Execute(v *Visit) error {
	return v.sendCustomEvent(ce.actionCause, ce.actionType, ce.customData)
}
