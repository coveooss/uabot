// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import (
	"errors"
	"fmt"
	"math/rand"

	ua "github.com/coveo/go-coveo/analytics"
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
	ignoreEvent   bool
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

	if se.ignoreEvent, validCast = e.Arguments["ignoreEvent"].(bool); !validCast {
		se.ignoreEvent = false
	}

	if se.goodQuery, validCast = e.Arguments["goodQuery"].(bool); !validCast {
		return nil, errors.New("Parameter goodQuery must be of type bool in SearchEvent")
	}
	if e.Arguments["customData"] != nil {
		if se.customData, validCast = e.Arguments["customData"].(map[string]interface{}); !validCast {
			return nil, errors.New("Parameter customData must be a json object (map[string]interface{}) in a search event")
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
// the analytics. Returns an error if something went wrong.
func (se *SearchEvent) Execute(v *Visit) (err error) {

	if se.query == "" { // if the query is empty, randomize one
		var queriesToRandom []string
		if queriesToRandom, err = se.getQueriesToRandomize(v); err != nil { // Figure out from which queries to randomize
			return
		}
		if se.query, err = randomQuery(queriesToRandom); err != nil { // Randomize the query from the selected array
			return
		}
	}
	se.keyword = se.query
	v.LastQuery.Q = se.query
	Info.Printf("Searching for : %s", se.query)

	// Execute a search and save the response
	if v.LastResponse, err = v.SearchClient.Query(*v.LastQuery); err != nil {
		return err
	}

	// in some scenarios (logging of page views), we don't want to send the search event to the analytics
	if !se.ignoreEvent {
		return se.send(v)
	}

	Info.Println("Ignoring the search event because of configuration.")
	return nil
}

func (se *SearchEvent) send(v *Visit) error {
	if v.LastResponse == nil {
		return errors.New("LastResponse was nil. Cannot send search event")
	}
	Info.Printf("Sending Search Event with %v results", v.LastResponse.TotalCount)
	event := ua.NewSearchEvent()

	v.DecorateEvent(event.ActionEvent)

	event.SearchQueryUID = v.LastResponse.SearchUID
	event.QueryText = se.query
	event.AdvancedQuery = v.LastQuery.AQ
	event.ActionCause = se.actionCause
	event.NumberOfResults = v.LastResponse.TotalCount
	event.ResponseTime = v.LastResponse.Duration

	v.DecorateCustomMetadata(event.ActionEvent, se.customData)

	if v.LastResponse.TotalCount > 0 {
		if urihash, ok := v.LastResponse.Results[0].Raw["sysurihash"].(string); ok {
			event.Results = []ua.ResultHash{
				ua.ResultHash{DocumentURI: v.LastResponse.Results[0].URI, DocumentURIHash: urihash},
			}
		} else {
			return errors.New("Cannot convert sysurihash to string in search event")
		}
	}

	// Send a UA search event
	return v.SendSearchEvent(event)
}

// getQueriesToRandomize Return an array of queries to randomize from.
func (se *SearchEvent) getQueriesToRandomize(v *Visit) (queriesToRandom []string, err error) {
	if se.goodQuery { // if we want a good query
		queriesToRandom = v.Config.GoodQueries
		if se.matchLanguage { // if the query must match the language
			if _, ok := v.Config.GoodQueriesInLang[v.Language]; !ok {
				err = errors.New("No good query detected in " + v.Language)
				return
			}
			queriesToRandom = v.Config.GoodQueriesInLang[v.Language]
		}
	} else { // if we want a bad query
		queriesToRandom = v.Config.BadQueries
		if se.matchLanguage { // if the query must match the language
			if _, ok := v.Config.BadQueriesInLang[v.Language]; !ok {
				err = errors.New("No bad query detected in " + v.Language)
				return
			}
			queriesToRandom = v.Config.BadQueriesInLang[v.Language]
		}
	}
	return
}

// randomQuery Returns a random query good or bad from the list of possible queries.
// returns an error if there are no queries to select from
func randomQuery(queries []string) (query string, err error) {
	if len(queries) < 1 {
		err = errors.New("Queries are empty")
		return
	}

	query = queries[rand.Intn(len(queries))]
	return
}
