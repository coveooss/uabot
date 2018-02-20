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
	clickRank     int
	offset        int
	probability   float64
	contentType   string
	pageViewField string
	customData    map[string]interface{}
}

func newViewEvent(e *JSONEvent, c *Config) (*ViewEvent, error) {
	var validcast bool
	var offset, docNo float64

	event := new(ViewEvent)

	if offset, validcast = e.Arguments["offset"].(float64); !validcast {
		return nil, errors.New("Parameter offset must be a positive number in a ViewEvent")
	}
	event.offset = int(offset)

	if event.probability, validcast = e.Arguments["probability"].(float64); !validcast || event.probability > 1 || event.probability < 0 {
		return nil, errors.New("Parameter probability must be a number between 0 and 1 in a ViewEvent")
	}

	if e.Arguments["contentType"] != nil {
		if event.contentType, validcast = e.Arguments["contentType"].(string); !validcast {
			return nil, errors.New("Parameter contentType must be of type string in ViewEvent")
		}
	} else {
		event.contentType = "Default ContentType"
	}

	if docNo, validcast = e.Arguments["docNo"].(float64); !validcast {
		return nil, errors.New("Parameter docNo must be a number in a ViewEvent")
	}
	event.clickRank = int(docNo)

	if e.Arguments["pageViewField"] != nil {
		if event.pageViewField, validcast = e.Arguments["pageViewField"].(string); !validcast {
			return nil, errors.New("Parameter pageViewField must be of type string in ViewEvent")
		}
	} else {
		event.pageViewField = c.RandomData.DefaultPageViewField
	}

	if e.Arguments["customData"] != nil {
		if event.customData, validcast = e.Arguments["customData"].(map[string]interface{}); !validcast {
			return nil, errors.New("Parameter customData must be a json object (map[string]interface{}) in a view event")
		}
	}

	return event, nil
}

// Execute the view event, sending a view event to the usage analytics
func (ve *ViewEvent) Execute(v *Visit) error {
	if v.LastResponse == nil {
		return errors.New("No query before pageView event, use a search event first")
	}
	if v.LastResponse.TotalCount < 1 {
		Warning.Printf("Last query %s returned no results cannot send view event", v.LastQuery.Q)
		return nil
	}

	if rand.Float64() <= ve.probability { // test if the event will exectute according to probability
		ve.clickRank = computeClickRank(v, ve.clickRank, ve.offset)

		if ve.clickRank > v.LastResponse.TotalCount {
			Warning.Printf("PageView index out of bounds, not sending event")
			return nil
		}

		return ve.send(v)
	}
	Info.Printf("User chose not to view (probability %v%%)", int(ve.probability*100))
	return nil
}

func (ve *ViewEvent) send(v *Visit) error {
	Info.Printf("Sending ViewEvent rank=%d ", ve.clickRank+1)

	event := ua.NewViewEvent()
	event.PageURI = v.LastResponse.Results[ve.clickRank].ClickURI
	event.PageTitle = v.LastResponse.Results[ve.clickRank].Title
	event.ContentType = ve.contentType
	event.ContentIDKey = "@" + ve.pageViewField
	event.PageReferrer = v.Referrer
	v.DecorateEvent(event.ActionEvent)
	v.DecorateCustomMetadata(event.ActionEvent, ve.customData)

	if _, ok := v.LastResponse.Results[ve.clickRank].Raw[ve.pageViewField]; !ok { // If the field does not exist on the "clicked" result
		Warning.Printf("Fields %s does not exist on result ranked %d. Not sending view event.", ve.pageViewField, ve.clickRank)
		return nil
	}
	if contentIDValue, ok := v.LastResponse.Results[ve.clickRank].Raw[ve.pageViewField].(string); ok { // If we can convert the fieldValue to a string
		event.ContentIDValue = contentIDValue
	} else {
		return fmt.Errorf("Cannot convert %s field %s value to string", v.LastResponse.Results[ve.clickRank].Raw[ve.pageViewField], ve.pageViewField)
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
