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
}

func newClickEvent(e *JSONEvent) (*ClickEvent, error) {
	var ok4, quickview bool
	offset, ok1 := e.Arguments["offset"].(float64)
	prob, ok2 := e.Arguments["probability"].(float64)
	rank, ok3 := e.Arguments["docNo"].(float64)
	if e.Arguments["quickview"] == nil {
		quickview = false
		ok4 = true
	} else {
		quickview, ok4 = e.Arguments["quickview"].(bool)
	}

	if !ok1 || !ok2 || !ok3 || !ok4 {
		return nil, errors.New("Invalid parse of arguments on Click Event")
	}

	if prob < 0 || prob > 1 {
		return nil, errors.New("Probability is out of bounds [0..1]")
	}

	return &ClickEvent{
		clickRank:   int(rank),
		offset:      int(offset),
		probability: prob,
		quickview:   quickview,
	}, nil
}

// Execute Execute the click event, sending a click event to the usage analytics
func (ce *ClickEvent) Execute(v *Visit) error {
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

		err := v.sendClickEvent(ce.clickRank, ce.quickview)
		if err != nil {
			return err
		}
		return nil
	}
	Info.Printf("User chose not to click (probability %v%%)", int(ce.probability*100))
	return nil
}
