package refactor

import (
	"encoding/json"
	"errors"
)

// ============== CLICK EVENT ======================
// ==================================================

// A ClickEvent is an event sent when the user clicks on a document
type ClickEvent struct {
	DocNo       int                    `json:"docNo,omitempty"`
	Offset      int                    `json:"offset,omitempty"`
	Probability float64                `json:"probability"`
	Quickview   bool                   `json:"quickview,omitempty"`
	CustomData  map[string]interface{} `json:"customData,omitempty"`
}

// Parse the remaining bits of the json event into the right arguments for this event.
func (e *ClickEvent) Parse(jse *JSONEvent) error {
	if err := json.Unmarshal(jse.Arguments, e); err != nil {
		return err
	}
	if e.Probability < 0 || e.Probability > 1 {
		return errors.New("Probability must be between 0 and 1 in a click event")
	}
	if e.Offset < 0 {
		return errors.New("Offset must be positive in a click event")
	}
	if e.DocNo < -1 {
		return errors.New("DocNo must be > 0 or -1 (for a random rank) in a click event")
	}
	return nil
}

// Execute Execute the search event, runs the query and sends a search event to
// the analytics.
func (e *ClickEvent) Execute(v *Visit) error {
	// Execute the event and send to analytics
	return nil
}

func (e *ClickEvent) sendTemplateEvent(v *Visit) error {
	// Send the actual analytics event
	return nil
}

// Check for interface implementation
var _ Event = (*ClickEvent)(nil)
