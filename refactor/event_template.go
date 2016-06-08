package refactor

import "encoding/json"

// ============== TEMPLATE EVENT ======================
// ==================================================

// A TemplateEvent is a template of an event to easily create a new event by copying this file.
type TemplateEvent struct {
	Argument1  string                 `json:"argument1,omitempty"`
	CustomData map[string]interface{} `json:"customData,omitempty"`
}

// Parse the remaining bits of the json event into the right arguments for this event.
func (e *TemplateEvent) Parse(jse *JSONEvent) error {
	if err := json.Unmarshal(jse.Arguments, e); err != nil {
		return err
	}
	return nil
}

// Execute Execute the search event, runs the query and sends a search event to
// the analytics.
func (e *TemplateEvent) Execute(v *Visit) error {
	// Execute the event and send to analytics
	return nil
}

func (e *TemplateEvent) sendTemplateEvent(v *Visit) error {
	// Send the actual analytics event
	return nil
}

// Check for interface implementation
var _ Event = (*TemplateEvent)(nil)
