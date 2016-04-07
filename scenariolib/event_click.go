// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import (
	"errors"
	"math"
	"math/rand"

	"github.com/k0kubun/pp"
)

// ============== CLICK EVENT ======================
// =================================================

// ClickEvent a struct representing a click, it is definied by a clickRank, an
// offset and a probability to click.
type ClickEvent struct {
	clickRank   int
	offset      int
	probability float64
}

func newClickEvent(e *JSONEvent) (*ClickEvent, error) {
	offset, ok1 := e.Arguments["offset"].(float64)
	prob, ok2 := e.Arguments["probability"].(float64)
	rank, ok3 := e.Arguments["docNo"].(float64)
	if !ok1 || !ok2 || !ok3 {
		return nil, errors.New("ERR >>> Invalid parse of arguments on Click Event")
	}

	if prob < 0 || prob > 1 {
		return nil, errors.New("ERR >>> Probability is out of bounds [0..1]")
	}

	return &ClickEvent{
		clickRank:   int(rank),
		offset:      int(offset),
		probability: prob,
	}, nil
}

// Execute Execute the click event, sending a click event to the usage analytics
func (ce *ClickEvent) Execute(v *Visit) error {
	if v.LastResponse.TotalCount < 1 {
		pp.Printf("WARN >>> Last query : %v returned no results, cannot click", v.LastQuery.Q)
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
			return errors.New("ERR >>> Click index out of bounds")
		}

		pp.Printf("\nLOG >>> Sending Click Event : Rank -> %v", ce.clickRank+1)

		err := v.sendClickEvent(ce.clickRank)
		if err != nil {
			return err
		}
		return nil
	}

	pp.Printf("\nLOG >>> User chose not to click with a probability of : %v %%", int(ce.probability*100))
	return nil
}
