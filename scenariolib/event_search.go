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
	query string
	// keyword exists because the query sent to the index may be different than the keyword(s) used to search
	keyword       string
	actionCause   string
	actionType    string
	logEvent      bool
	customData    map[string]interface{}
	matchLanguage bool
	goodQuery     bool
}

func newSearchEvent(e *JSONEvent, c *Config) (*SearchEvent, error) {
	var inputTitle string
	var validCast bool
	se := new(SearchEvent)

	if se.query, validCast = e.Arguments["queryText"].(string); !validCast {
		return nil, errors.New("Parameter query must be of type string in SearchEvent")
	}
	if e.Arguments["logEvent"] != nil {
		if se.logEvent, validCast = e.Arguments["logEvent"].(bool); !validCast {
			return nil, errors.New("Parameter logEvent must be of type bool in SearchEvent")
		}
	} else {
		se.logEvent = true
	}
	Info.Printf("Will log search event to analytics: (%t)", se.logEvent)

	if se.goodQuery, validCast = e.Arguments["goodQuery"].(bool); !validCast {
		return nil, errors.New("Parameter goodQuery must be of type bool in SearchEvent")
	}
	if e.Arguments["customData"] != nil {
		if se.customData, validCast = e.Arguments["customData"].(map[string]interface{}); !validCast {
			return nil, errors.New("Parameter custom must be a json object (map[string]interface{}) in a search event.")
		}
	}

	if e.Arguments["matchLanguage"] != nil {
		if se.matchLanguage, validCast = e.Arguments["matchLanguage"].(bool); !validCast {
			return nil, errors.New("Parameter matchLanguage must be a type bool in SearchEvent")
		}
	}

	se.keyword = se.query
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
			se.query = fmt.Sprintf("($some(keywords: %s, match: 1, removeStopWords: true, maximum: 300)) ($sort(criteria: relevancy))", se.keyword)
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
	var err error
	if se.query == "" {
		if se.matchLanguage {
			se.query, err = v.Config.RandomQueryInLanguage(se.goodQuery, v.Language)
			se.keyword = se.query
			if err != nil {
				return err
			}
		} else {
			se.query, err = v.Config.RandomQuery(se.goodQuery)
			se.keyword = se.query
			if err != nil {
				return err
			}
		}
	}

	Info.Printf("Searching for : %s", se.query)

	v.LastQuery.Q = se.query

	// Execute a search and save the response
	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}
	v.LastResponse = resp

	// in some scenarios (logging of page views), we don't want to send the search event to the analytics
	if se.logEvent {
		err = v.sendSearchEvent(se.keyword, se.actionCause, se.actionType, se.customData)
		if err != nil {
			return err
		}
	}
	return nil
}
