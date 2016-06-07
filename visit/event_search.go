package visit

import "encoding/json"

// Search Contains the arguments necessary to send a search event to the analytics
type Search struct {
	QueryText  string `json:"queryText"`
	GoodQuery  bool   `json:"goodQuery"`
	CaseSearch bool   `json:"caseSearch,omitempty"`
	InputTitle string `json:"inputTitle,omitempty"`
	// keyword exists because the query sent to the index may be diffrent than the keyword(s) used to search
	keyword     string
	actionCause string
	actionType  string
	customData  map[string]interface{}
}

// Parse Parse the different arguments in the JSONEvent to build the event
func (e *Search) Parse(jse *JSONEvent) error {
	err := json.Unmarshal(jse.Arguments, e)
	if err != nil {
		return err
	}
	return nil
}

// Execute Send the event to the analytics endpoint
func (e *Search) Execute(v *Visit) error {
	return nil
}

// Check for interface implementation
var _ Event = (*Search)(nil)
