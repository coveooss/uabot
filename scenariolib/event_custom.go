// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

// ============== SEARCH EVENT ======================
// ==================================================

// CustomEvent a struct representing a search, is defined by a query to execute
type CustomEvent struct {
	EventType  string                 `json:"eventType"`
	EventValue string                 `json:"eventValue"`
	CustomData map[string]interface{} `json:"customData,omitempty"`
}

// IsValid Additional validation after the json unmarshal.
func (custom *CustomEvent) IsValid() (bool, string) {
	return true, ""
}

// Execute the search event, runs the query and sends a search event to
// the analytics.
func (custom *CustomEvent) Execute(v *Visit) error {
	return v.sendCustomEvent(custom.EventValue, custom.EventType, custom.CustomData)
}
