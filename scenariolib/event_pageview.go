package scenariolib

import (
	"errors"
	"math"
	"math/rand"
)

// ============== VIEW EVENT ======================
// =================================================

// ViewEvent a struct representing a view, it is defined by a clickRank, an
// offset, a probability to click, the contentType and a pageViewField
// We keep a similar structure to a click event because we can simulate that the user
// visited a page that was returned as a result
type ViewEvent struct {
	clickRank     int
	offset        int
	probability   float64
	contentType   string
	pageViewField string
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
		event.pageViewField = c.DefaultPageViewField
	}

	return event, nil
}

// Execute the view event, sending a view event to the usage analytics
func (ve *ViewEvent) Execute(v *Visit) error {
	if v.LastResponse == nil {
		return errors.New("LastResponse was nil, cannot send a pageview. Please use a search event before.")
	}
	if v.LastResponse.TotalCount < 1 {
		Warning.Printf("Last query %s returned no results cannot click", v.LastQuery.Q)
		return nil
	}

	if rand.Float64() <= ve.probability { // test if the event will exectute according to probability
		ve.computeClickRank(v)

		err := v.sendViewEvent(ve.clickRank, ve.contentType, ve.pageViewField)
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
