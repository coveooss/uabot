package visit

import "encoding/json"

// Scenario A definition of a single visit to the search page.
// Contains events to define the actions in the visit.
type Scenario struct {
	Name      string      `json:"name"`
	Weight    int         `json:"weight"`
	UserAgent string      `json:"useragent,omitempty"`
	Events    []JSONEvent `json:"events"`
}

// ScenarioMap Type used to create a weighted random structure of scenarios.
type ScenarioMap map[int]*Scenario

// JSONEvent An action taken by the user such as a search, a click, a SearchAndClick, etc.
type JSONEvent struct {
	Type      string          `json:"type"`
	Arguments json.RawMessage `json:"arguments"`
}
