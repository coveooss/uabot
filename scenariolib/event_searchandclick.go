// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import (
	"errors"
	"math/rand"

	"github.com/k0kubun/pp"
)

// ============== SEARCH AND CLICK EVENT ======================
// ============================================================

// SearchAndClickEvent represents a search event followed by a click on a specific
// document found by the title
type SearchAndClickEvent struct {
	query    string
	docTitle string
	prob     float64
}

func newSearchAndClickEvent(e *JSONEvent) (*SearchAndClickEvent, error) {
	query, ok1 := e.Arguments["queryText"].(string)
	docTitle, ok2 := e.Arguments["docClickTitle"].(string)
	prob, ok3 := e.Arguments["probability"].(float64)
	if !ok1 || !ok2 || !ok3 {
		return nil, errors.New("ERR >>> Invalid parse of arguments on SearchAndClick Event")
	}

	return &SearchAndClickEvent{
		query:    query,
		docTitle: docTitle,
		prob:     prob,
	}, nil
}

// Execute Execute the search and click event sending both events to the analytics
func (sc *SearchAndClickEvent) Execute(v *Visit) error {
	// Execute the search event
	se := new(SearchEvent)
	se.query = sc.query
	err := se.Execute(v)
	if err != nil {
		return err
	}

	if v.LastResponse.TotalCount < 1 {
		return errors.New("ERR >>> Last query returned no results")
	}

	WaitBetweenActions()

	if rand.Float64() <= sc.prob {
		rank := v.FindDocumentRankByTitle(sc.docTitle)
		if rank >= 0 {
			pp.Printf("\nLOG >>> Sending Click Event => Found Document Rank: %v", rank+1)

			ce := new(ClickEvent)
			ce.clickRank = rank
			ce.offset = 0
			ce.probability = 1

			ce.Execute(v)
			if err != nil {
				return err
			}
		} else {
			return errors.New("ERR >>> Could not find the specific document you are looking for")
		}
	} else {
		pp.Printf("\nLOG >>> User chose not to click with a probability of : %v %%", int(sc.prob*100))
	}

	return nil
}
