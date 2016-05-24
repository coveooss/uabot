// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import (
	"errors"
	"fmt"
)

// ============== SEARCH EVENT ======================
// ==================================================

// SearchEvent a struct representing a search, is defined by a query to execute
type SearchEvent struct {
	query       string
	actionCause string
	actionType  string
	customData  map[string]interface{}
}

func newSearchEvent(e *JSONEvent, c *Config) (*SearchEvent, error) {
	var err error
	var inputTitle string
	var goodQuery, validCast bool

	se := new(SearchEvent)

	if se.query, validCast = e.Arguments["queryText"].(string); !validCast {
		return nil, errors.New("Parameter query must be of type string in SearchEvent")
	}
	if goodQuery, validCast = e.Arguments["goodQuery"].(bool); !validCast {
		return nil, errors.New("Parameter goodQuery must be of type bool in SearchEvent")
	}

	if se.query == "" {
		se.query, err = c.RandomQuery(goodQuery)
		if err != nil {
			return nil, err
		}
	}

	se.actionCause = "searchboxSubmit"
	se.actionType = "search box"

	if e.Arguments["caseSearch"] != nil {
		caseSearch, validCast := e.Arguments["caseSearch"].(bool)
		if !validCast {
			return nil, errors.New("Parameter caseSearch must be a boolean")
		}
		if caseSearch {
			se.actionCause = "inputChange"
			se.actionType = "caseCreation"
			keyword := se.query
			se.query = fmt.Sprintf("($some(keywords: %s, match: 1, removeStopWords: true, maximum: 300)) ($qre(expression: undefined=%s, modifier: 50))", keyword, keyword)
			if inputTitle, validCast = e.Arguments["inputTitle"].(string); !validCast {
				return nil, errors.New("Parameter inputTitle is required in a caseSearch and must be a string")
			}
			se.customData = map[string]interface{}{
				"inputTitle": inputTitle,
			}
		}
	}

	return se, nil
}

// Execute Execute the search event, runs the query and sends a search event to
// the analytics.
func (se *SearchEvent) Execute(v *Visit) error {
	Info.Printf("Searching for : %s", se.query)
	v.LastQuery.Q = se.query

	// Execute a search and save the response
	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}
	v.LastResponse = resp

	err = v.sendSearchEvent(se.query, se.actionCause, se.actionType, se.customData)
	if err != nil {
		return err
	}
	return nil
}
