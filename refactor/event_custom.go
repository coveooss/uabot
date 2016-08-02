package refactor

import (
	"encoding/json"
	"math/rand"

	"github.com/go-coveo/analytics"
)

// ============== CUSTOM EVENT ======================
// ==================================================

// A CustomEvent is a custom event sent to the analytics
type CustomEvent struct {
	EventType  string                 `json:"eventType"`
	EventValue string                 `json:"eventValue"`
	CustomData map[string]interface{} `json:"customData,omitempty"`
}

// Parse the remaining bits of the json event into the right arguments for this event.
func (e *CustomEvent) Parse(jse *JSONEvent) error {
	if err := json.Unmarshal(jse.Arguments, e); err != nil {
		return err
	}
	return nil
}

// Execute the custom event, runs the query and sends a search event to
// the analytics.
func (e *CustomEvent) Execute(v *Visit) error {

	defaultCustomData := map[string]interface{}{
		"JSUIVersion": JSUIVERSION,
		"ipadress":    v.User.IP,
	}

	if e.CustomData == nil {
		e.CustomData = defaultCustomData
	} else {
		for k, v := range defaultCustomData {
			if e.CustomData[k] == nil {
				e.CustomData[k] = v
			}
		}
	}

	// Execute the event and send to analytics
	if err := e.sendCustomEvent(v); err != nil {
		return err
	}
	return nil
}

func (e *CustomEvent) sendCustomEvent(v *Visit) error {
	// Info.Printf("Sending CustomEvent cause: %s ||| type: %s", actionCause, actionType)
	event := analytics.NewCustomEvent()

	event.Username = v.User.Email
	event.Anonymous = v.User.Anonymous
	event.OriginLevel1 = v.OriginLevel1
	event.OriginLevel2 = v.OriginLevel2
	event.Language = v.User.Language
	event.EventType = e.EventType
	event.EventValue = e.EventValue
	event.CustomData = e.CustomData

	// Send all the possible random custom data that can be added from the config
	// scenario file.
	for _, elem := range v.Config.RandomCustomData {
		event.CustomData[elem.APIName] = elem.Values[rand.Intn(len(elem.Values))]
	}

	// Send the actual analytics event
	err := v.UAClient.SendCustomEvent(*event)
	return err
}

// Check for interface implementation
var _ Event = (*CustomEvent)(nil)
