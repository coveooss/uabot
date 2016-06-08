package refactor

import "encoding/json"

// JSONEvent An action taken by the user such as a search, a click, a SearchAndClick, etc.
type JSONEvent struct {
	Type      string          `json:"type"`
	Arguments json.RawMessage `json:"arguments"`
}
