// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import (
	"errors"
	"math/rand"
)

// ============== SEARCH AND CLICK EVENT ======================
// ============================================================

// SearchAndClickEvent represents a search event followed by a click on a specific
// document found by the title
type SearchAndClickEvent struct {
	query     string
	docTitle  string
	prob      float64
	quickview bool
}

func newSearchAndClickEvent(e *JSONEvent) (*SearchAndClickEvent, error) {
	var quickview, ok4 bool
	query, ok1 := e.Arguments["queryText"].(string)
	docTitle, ok2 := e.Arguments["docClickTitle"].(string)
	prob, ok3 := e.Arguments["probability"].(float64)
	if e.Arguments["quickview"] == nil {
		quickview = false
		ok4 = true
	} else {
		quickview, ok4 = e.Arguments["quickview"].(bool)
	}
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return nil, errors.New("Invalid parse of arguments on SearchAndClick Event")
	}

	return &SearchAndClickEvent{
		query:     query,
		docTitle:  docTitle,
		prob:      prob,
		quickview: quickview,
	}, nil
}

// Execute Execute the search and click event sending both events to the analytics
func (sc *SearchAndClickEvent) Execute(v *Visit) error {
	Info.Printf("Executing SearchAndClickEvent : Searching for %s, clicking on %s (quickview %v)", sc.query, sc.docTitle, sc.quickview)
	// Execute the search event
	se := new(SearchEvent)
	se.query = sc.query
	err := se.Execute(v)
	if err != nil {
		return err
	}

	if v.LastResponse.TotalCount < 1 {
		return errors.New("Last query returned no results")
	}

	WaitBetweenActions()

	if rand.Float64() <= sc.prob {
		rank := v.FindDocumentRankByTitle(sc.docTitle)
		if rank >= 0 {
			Info.Printf("Sending ClickEvent => Found document at rank : %d", rank+1)

			ce := new(ClickEvent)
			ce.clickRank = rank
			ce.offset = 0
			ce.probability = 1

			ce.quickview = sc.quickview

			ce.Execute(v)
			if err != nil {
				return err
			}
		} else {
			return errors.New("Could not find the specific document you are looking for")
		}
	} else {
		Info.Printf("User chose not to click (probability %v%%)", int(sc.prob*100))
	}

	return nil
}
