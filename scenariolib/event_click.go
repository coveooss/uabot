package scenariolib

import (
	"errors"
	"math"
	"math/rand"
)

// ============== CLICK EVENT ======================
// =================================================

// ClickEvent a struct representing a click, it is definied by a clickRank, an
// offset and a probability to click.
type ClickEvent struct {
	clickRank   int
	offset      int
	probability float64
	quickview   bool
	customData  map[string]interface{}
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

	if docNo, validcast = e.Arguments["offset"].(float64); !validcast {
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

	return event, nil
}

// Execute Execute the click event, sending a click event to the usage analytics
func (ce *ClickEvent) Execute(v *Visit) error {
	if ce.clickRank == -1 { // if rank == -1 we need to randomize a rank
		ce.clickRank = 0
		// Find a random rank within the possible click values accounting for the offset
		if v.LastResponse.TotalCount > 1 {
			topL := Min(v.LastQuery.NumberOfResults, v.LastResponse.TotalCount)
			rndRank := int(math.Abs(rand.NormFloat64()*2)) + ce.offset
			ce.clickRank = Min(rndRank, topL-1)
		} else {
			ce.clickRank = 1
		}
	}

	if rand.Float64() <= ce.probability { // Probability to click
		if ce.clickRank > v.LastResponse.TotalCount {
			return errors.New("Click index out of bounds")
		}

		err := v.sendClickEvent(ce.clickRank, ce.quickview, ce.customData, v.LastResponse.TotalCount > 0)
		if err != nil {
			return err
		}
		return nil
	}
	Info.Printf("User chose not to click (probability %v%%)", int(ce.probability*100))
	return nil
}
