package scenariolib

import (
	"encoding/json"
)

// Scenario Represents one visit to the search
type Scenario struct {
	// A Name given to the scenario for easier logging.
	Name string `json:"name"`

	// A Weight for randomizing scenarios.
	Weight int `json:"weight"`

	// A UserAgent string representing the visit
	UserAgent string `json:"useragent,omitempty"`

	// An Events array of actions the user will take
	Events []JSONEvent `json:"events"`

	// A Language for this scenario
	Language string `json:"language,omitempty"`

	// Mobile A boolean value if this visit is forced on mobile
	Mobile bool `json:"mobile,omitempty"`
}

// JSONEvent An action taken by the user such as a search, a click, a SearchAndClick, etc.
// Type A string describing the type of event
// Arguments An array of the arguments to the event, specific to the type of event.
type JSONEvent struct {
	Type      string          `json:"type"`
	Arguments json.RawMessage `json:"arguments"`
}
