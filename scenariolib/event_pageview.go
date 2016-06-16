package scenariolib

import (
	"errors"
	"math"
	"math/rand"
)

// ============== VIEW EVENT ======================
// =================================================

// ViewEvent a struct representing a view, it is defined by a clickRank, an
// offset, a probability to click and the contentType
type ViewEvent struct {
	clickRank   int
	offset      int
	probability float64
	contentType string
}

func newViewEvent(e *JSONEvent) (*ViewEvent, error) {
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

	return event, nil
}

// Execute Execute the view event, sending a view event to the usage analytics
func (ve *ViewEvent) Execute(v *Visit) error {
	if v.LastResponse.TotalCount < 1 {
		Warning.Printf("Last query %s returned no results cannot click", v.LastQuery.Q)
		return nil
	}

	if ve.clickRank == -1 { // if rank == -1 we need to randomize a rank
		ve.clickRank = 0
		// Find a random rank within the possible click values accounting for the offset
		if v.LastResponse.TotalCount > 1 {
			topL := Min(v.LastQuery.NumberOfResults, v.LastResponse.TotalCount)
			rndRank := int(math.Abs(rand.NormFloat64()*2)) + ve.offset
			ve.clickRank = Min(rndRank, topL-1)
		}
	}

	if rand.Float64() <= ve.probability { // Probability to click
		if ve.clickRank > v.LastResponse.TotalCount {
			return errors.New("Search results index out of bounds")
		}

		err := v.sendViewEvent(ve.clickRank, ve.contentType)
		if err != nil {
			return err
		}
		return nil
	}
	Info.Printf("User chose not to view (probability %v%%)", int(ve.probability*100))
	return nil
}
