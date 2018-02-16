package scenariolib

import (
	"errors"
	"math"
	"math/rand"
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
		ve.computeClickRank(v)

		if ve.clickRank > v.LastResponse.TotalCount {
			Warning.Printf("PageView index out of bounds, not sending event")
			return nil
		}

		err := v.sendViewEvent(ve.clickRank, ve.contentType, ve.pageViewField, ve.customData)
		if err != nil {
			return err
		}
		return nil
	}
	Info.Printf("User chose not to view (probability %v%%)", int(ve.probability*100))
	return nil
}

// Randomize a click rank if the clickRank is -1
func (ve *ViewEvent) computeClickRank(v *Visit) *ViewEvent {
	if ve.clickRank == -1 { // if rank == -1 we need to randomize a rank
		ve.clickRank = 0
		// Find a random rank within the possible click values accounting for the offset
		if v.LastResponse.TotalCount > 1 {
			topL := Min(v.LastQuery.NumberOfResults, v.LastResponse.TotalCount)
			rndRank := int(math.Abs(rand.NormFloat64()*2)) + ve.offset
			ve.clickRank = Min(rndRank, topL-1)
		}
	}
	return ve
}
