package scenariolib

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	ua "github.com/coveo/go-coveo/analytics"
)

// ============== VIEW EVENT ======================
// =================================================

// ViewEvent a struct representing a view, it is defined by a clickRank, an offset, a probability to click,
// the contentType and a pageViewField. We keep a similar structure to a click event because we can simulate
// that the user visited a page that was returned as a result from a search.
// The view event sent will contain {contentIdKey: "@pageViewField", contentIdValue: result['pageViewField']}
type ViewEvent struct {
	ClickRank     int                    `json:"clickRank"`
	Probability   float64                `json:"probability"`
	PageViewField string                 `json:"pageViewField"`
	Offset        int                    `json:"offset,omitempty"`
	ContentType   string                 `json:"contentType,omitempty"`
	CustomData    map[string]interface{} `json:"customData,omitempty"`
}

// IsValid Additional validation after the json unmarshal.
func (view *ViewEvent) IsValid() (bool, string) {
	if view.Probability < 0 || view.Probability > 1 {
		return false, "A view event probability must be between 0 and 1."
	}
	return true, ""
}

// Execute the view event, sending a view event to the usage analytics
func (view *ViewEvent) Execute(v *Visit) error {
	if v.LastResponse == nil {
		return errors.New("No query before pageView event, use a search event first")
	}
	if v.LastResponse.TotalCount < 1 {
		Warning.Printf("Last query %s returned no results cannot send view event", v.LastQuery.Q)
		return nil
	}

	if rand.Float64() <= view.Probability { // test if the event will exectute according to probability
		view.ClickRank = computeClickRank(v, view.ClickRank, view.Offset)

		if view.ClickRank > v.LastResponse.TotalCount {
			Warning.Printf("PageView index out of bounds, not sending event")
			return nil
		}

		return view.send(v)
	}
	Info.Printf("User chose not to view (probability %v%%)", int(view.Probability*100))
	return nil
}

func (view *ViewEvent) send(v *Visit) error {
	Info.Printf("Sending ViewEvent rank=%d ", view.ClickRank+1)

	event := ua.NewViewEvent()
	event.Location = v.LastResponse.Results[view.ClickRank].ClickURI
	event.Title = v.LastResponse.Results[view.ClickRank].Title
	event.ContentType = view.ContentType
	event.ContentIDKey = "@" + view.PageViewField
	event.Referrer = v.Referrer
	v.DecorateEvent(event.ActionEvent)
	v.DecorateCustomMetadata(event.ActionEvent, view.CustomData)

	if _, ok := v.LastResponse.Results[view.ClickRank].Raw[view.PageViewField]; !ok { // If the field does not exist on the "clicked" result
		Warning.Printf("Fields %s does not exist on result ranked %d. Not sending view event.", view.PageViewField, view.ClickRank)
		return nil
	}
	if contentIDValue, ok := v.LastResponse.Results[view.ClickRank].Raw[view.PageViewField].(string); ok { // If we can convert the fieldValue to a string
		event.ContentIDValue = contentIDValue
	} else {
		return fmt.Errorf("Cannot convert %s field %s value to string", v.LastResponse.Results[view.ClickRank].Raw[view.PageViewField], view.PageViewField)
	}

	// Send a UA view event
	return v.SendViewEvent(event)
}

// Randomize a click rank if the clickRank is -1
func computeClickRank(v *Visit, setRank, offset int) (clickRank int) {
	if setRank == -1 { // if rank == -1 we need to randomize a rank
		// Find a random rank within the possible click values accounting for the offset
		if v.LastResponse.TotalCount > 1 {
			topL := Min(v.LastQuery.NumberOfResults, v.LastResponse.TotalCount)
			rndRank := int(math.Abs(rand.NormFloat64()*2)) + offset
			clickRank = Min(rndRank, topL-1)
		}
	}
	return
}
