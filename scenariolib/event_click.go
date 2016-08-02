package scenariolib

import (
	"errors"
	"math"
	"math/rand"

	"encoding/json"

	"github.com/coveo/go-coveo/search"
)

// ============== CLICK EVENT ======================
// =================================================

// ClickEvent a struct representing a click, it is defined by a clickRank, an
// offset and a probability to click.
type ClickEvent struct {
	clickRank    int
	offset       int
	probability  float64
	quickview    bool
	customData   map[string]interface{}
	fakeClick    bool
	fakeResponse search.Response
}

func newClickEvent(e *JSONEvent) (*ClickEvent, error) {
	var validcast bool
	var offset, docNo float64

	event := new(ClickEvent)

	if offset, validcast = e.Arguments["offset"].(float64); !validcast {
		return nil, errors.New("Parameter offset must be a positive number in a ClickEvent")
	}
	event.offset = int(offset)

	if event.probability, validcast = e.Arguments["probability"].(float64); !validcast || event.probability > 1 || event.probability < 0 {
		return nil, errors.New("Parameter probability must be a number between 0 and 1 in a ClickEvent")
	}

	if docNo, validcast = e.Arguments["docNo"].(float64); !validcast {
		return nil, errors.New("Parameter docNo must be a number in a ClickEvent")
	}
	event.clickRank = int(docNo)

	if e.Arguments["quickview"] != nil {
		if event.quickview, validcast = e.Arguments["quickview"].(bool); !validcast {
			return nil, errors.New("Parameter quickview must be a boolean")
		}
	} else {
		event.quickview = false
	}

	if e.Arguments["customData"] != nil {
		if event.customData, validcast = e.Arguments["customData"].(map[string]interface{}); !validcast {
			return nil, errors.New("Parameter custom must be a json object (map[string]interface{}) in a click event.")
		}
	}

	if e.Arguments["fakeClick"] != nil {
		if event.fakeClick, validcast = e.Arguments["fakeClick"].(bool); !validcast {
			return nil, errors.New("Parameter fakeClick must be a boolean value")
		}
		if e.Arguments["falseResponse"] != nil {
			jsonFalseResponse, _ := json.Marshal(e.Arguments["falseResponse"])
			err := json.Unmarshal(jsonFalseResponse, &event.fakeResponse)
			if err != nil {
				return nil, errors.New("Parameter falseResponse must be a search.Response")
			}
		}
	} else {
		event.fakeClick = false
	}
	return event, nil
}

// Execute Execute the click event, sending a click event to the usage analytics
func (ce *ClickEvent) Execute(v *Visit) error {
<<<<<<< HEAD
	if ce.fakeClick {
		searchUID := v.LastResponse.SearchUID
		v.LastResponse = &ce.fakeResponse
		v.LastResponse.SearchUID = searchUID
	}

    // Error handling, error if last response is nil, warning if last response had no results
	if v.LastResponse == nil {
		return errors.New("Cannot execute a click on a nil LastResponse. Please use a search event first.")
	}

	if v.LastResponse.TotalCount < 1 {
		Warning.Printf("Last query %s returned no results cannot click", v.LastQuery.Q)
		return nil
	}
	
	if ce.clickRank == -1 { // if rank == -1 we need to randomize a rank
		ce.clickRank = 0
		// Find a random rank within the possible click values accounting for the offset
		if v.LastResponse.TotalCount > 1 {
			topL := Min(v.LastQuery.NumberOfResults, v.LastResponse.TotalCount)
			rndRank := int(math.Abs(rand.NormFloat64()*2)) + ce.offset
			ce.clickRank = Min(rndRank, topL-1)
		}
	}

	if rand.Float64() <= ce.probability { // Probability to click
		if ce.clickRank > v.LastResponse.TotalCount {
			return errors.New("Click index out of bounds")
		}

		err := v.sendClickEvent(ce.clickRank, ce.quickview, ce.customData)
		if err != nil {
			return err
		}
		return nil
	}
	Info.Printf("User chose not to click (probability %v%%)", int(ce.probability*100))
	return nil
}
