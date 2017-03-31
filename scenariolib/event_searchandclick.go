package scenariolib

import (
	"errors"
	"fmt"
	"math/rand"
)

// ============== SEARCH AND CLICK EVENT ======================
// ============================================================

// SearchAndClickEvent represents a search event followed by a click on a specific
// document found by the title
type SearchAndClickEvent struct {
	query      string
	docTitle   string
	prob       float64
	quickview  bool
	caseSearch bool
	inputTitle string
	customData map[string]interface{}
}

func newSearchAndClickEvent(e *JSONEvent) (*SearchAndClickEvent, error) {
	var query, docClickTitle, inputTitle string
	var prob float64
	var quickview, caseSearch, validCast bool
	var customData map[string]interface{}

	if query, validCast = e.Arguments["queryText"].(string); !validCast {
		return nil, errors.New("Parameter queryText must be of type string in SearchAndClickEvent")
	}

	if docClickTitle, validCast = e.Arguments["docClickTitle"].(string); !validCast {
		return nil, errors.New("Parameter docClickTitle must be of type string in SearchAndClickEvent")
	}

	if prob, validCast = e.Arguments["probability"].(float64); !validCast {
		return nil, errors.New("Parameter probability must be of type float64 in SearchAndClickEvent")
	}

	if e.Arguments["quickview"] == nil {
		quickview = false
	} else {
		if quickview, validCast = e.Arguments["quickview"].(bool); !validCast {
			return nil, errors.New("Parameter quickview must be of type boolean in SearchAndClickEvent")
		}
	}

	if e.Arguments["caseSearch"] != nil {
		if caseSearch, validCast = e.Arguments["caseSearch"].(bool); !validCast {
			return nil, errors.New("Parameter caseSearch must be of type boolean in SearchAndClickEvent")
		}
		if caseSearch {
			if inputTitle, validCast = e.Arguments["inputTitle"].(string); !validCast {
				return nil, errors.New("Parameter inputTitle is mandatory on a caseSearch and must be of type string in SearchAndClickEvent")
			}
		}
	}

	if e.Arguments["customData"] != nil {
		if customData, validCast = e.Arguments["customData"].(map[string]interface{}); !validCast {
			return nil, errors.New("Parameter custom must be a json object (map[string]interface{}) in a custom event.")
		}
	}

	return &SearchAndClickEvent{
		query:      query,
		docTitle:   docClickTitle,
		prob:       prob,
		quickview:  quickview,
		caseSearch: caseSearch,
		inputTitle: inputTitle,
		customData: customData,
	}, nil
}

// Execute Execute the search and click event sending both events to the analytics
func (sc *SearchAndClickEvent) Execute(v *Visit) error {
	Info.Printf("Executing SearchAndClickEvent : Searching for %s, clicking on %s (quickview %v)", sc.query, sc.docTitle, sc.quickview)
	// Execute the search event
	se := new(SearchEvent)
	se.query = sc.query
	se.keyword = sc.query
	se.customData = make(map[string]interface{})
	if sc.caseSearch {
		se.query = fmt.Sprintf("($some(keywords: %s, match: 1, removeStopWords: true, maximum: 300)) ($sort(criteria: relevancy))", se.query)
		se.actionCause = "inputChange"
		se.actionType = "caseCreation"
		se.customData["inputTitle"] = sc.inputTitle
	} else {
		se.actionCause = "searchboxSubmit"
		se.actionType = "search box"
	}
	// Override possible values of customData with the specific customData sent
	for k, v := range sc.customData {
		se.customData[k] = v
	}
	err := se.Execute(v)
	if err != nil {
		return err
	}

	if v.LastResponse.TotalCount < 1 {
		return errors.New("Last query returned no results")
	}

	var timeToWait int
	if v.Config.TimeBetweenActions > 0 {
		timeToWait = v.Config.TimeBetweenActions
	} else {
		timeToWait = DEFAULTTIMEBETWEENACTIONS
	}
	WaitBetweenActions(timeToWait, v.Config.IsWaitConstant)

	if rand.Float64() <= sc.prob {
		rank := v.FindDocumentRankByTitle(sc.docTitle)
		if rank >= 0 {
			Info.Printf("Sending ClickEvent => Found document at rank : %d", rank+1)

			ce := new(ClickEvent)
			ce.clickRank = rank
			ce.offset = 0
			ce.probability = 1
			ce.quickview = sc.quickview

			ce.customData = make(map[string]interface{})
			// Override possible values of customData with the specific customData sent
			for k, v := range sc.customData {
				ce.customData[k] = v
			}
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
