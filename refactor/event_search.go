package refactor

import "encoding/json"

// ============== SEARCH EVENT ======================
// ==================================================

// A TemplateEvent is a template of an event to easily create a new event by copying this file.
type SearchEvent struct {
	QueryText     string                 `json:"queryText"`
	GoodQuery     bool                   `json:"goodQuery"`
	MatchLanguage bool                   `json:"matchLanguage,omitempty"`
	MockEvent     bool                   `json:"mockEvent,omitempty"`
	CaseSearch    bool                   `json:"caseSearch,omitempty"`
	InputTitle    string                 `json:"inputTitle,omitempty"`
	CustomData    map[string]interface{} `json:"customData,omitempty"`
	query         string
	actionCause   string
	actionType    string
}

// Parse the remaining bits of the json event into the right arguments for this event.
func (e *SearchEvent) Parse(jse *JSONEvent) error {
	if err := json.Unmarshal(jse.Arguments, e); err != nil {
		return err
	}
	// Info.Printf("Mocking event: (%t)", e.MockEvent)
	e.actionCause = "searchboxSubmit"
	return nil
}

// Execute the search event, runs the query and sends a search event to
// the analytics.
func (e *SearchEvent) Execute(v *Visit) error {
	// Execute the event and send to analytics
	return nil
}

func (e *SearchEvent) sendSearchEvent(v *Visit) error {
	// Send the actual analytics event
	return nil
}

// Check for interface implementation
var _ Event = (*SearchEvent)(nil)
