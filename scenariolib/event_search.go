// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import (
	"errors"

	"github.com/k0kubun/pp"
)

// ============== SEARCH EVENT ======================
// ==================================================

// SearchEvent a struct representing a search, is defined by a query to execute
type SearchEvent struct {
	query string
}

func newSearchEvent(e *JSONEvent, c *Config) (*SearchEvent, error) {
	var err error
	query, ok1 := e.Arguments["queryText"].(string)
	goodQuery, ok2 := e.Arguments["goodQuery"].(bool)
	if !ok1 || !ok2 {
		return nil, errors.New("ERR >>> Invalid parse of arguments on Search Event")
	}

	if query == "" {
		query, err = c.RandomQuery(goodQuery)
		if err != nil {
			return nil, err
		}
	}

	return &SearchEvent{
		query: query,
	}, nil
}

// Execute Execute the search event, runs the query and sends a search event to
// the analytics.
func (se *SearchEvent) Execute(v *Visit) error {
	pp.Printf("\nLOG >>> Searching for : %v", se.query)
	v.LastQuery.Q = se.query

	// Execute a search and save the response
	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}
	v.LastResponse = resp

	err = v.sendSearchEvent(se.query)
	if err != nil {
		return err
	}
	return nil
}
