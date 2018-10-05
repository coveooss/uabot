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
	Query         string                 `json:"queryText,omitempty"`
	IgnoreEvent   bool                   `json:"ignoreEvent,omitempty"`
	GoodQuery     bool                   `json:"goodQuery,omitempty"`
	ActionCause   string                 `json:"actionCause,omitempty"`
	CaseSearch    bool                   `json:"caseSearch,omitempty"`
	InputTitle    string                 `json:"inputTitle,omitempty"`
	MatchLanguage bool                   `json:"matchLanguage,omitempty"`
	CustomData    map[string]interface{} `json:"customData,omitempty"`
	Keyword       string
	ActionType    string
}

const caseQuerySomeTemplate = "($some(keywords: %s, match: 1, removeStopWords: true, maximum: 300)) ($sort(criteria: relevancy))"
const defaultSearchCause = "searchboxSubmit"
const defaultCaseSearchCause = "inputChange"

// IsValid Additional validation after the json unmarshal.
func (search *SearchEvent) IsValid() (bool, string) {
	if search.CaseSearch && search.InputTitle == "" {
		return false, "If caseSearch is true, you need to provide an inputTitle."
	}
	return true, ""
}

func (search *SearchEvent) handleCaseSearch() {
	Info.Println("Executing a Case Search.")
	search.ActionCause = defaultCaseSearchCause
	search.ActionType = "caseCreation"
	search.Query = fmt.Sprintf(caseQuerySomeTemplate, search.Keyword)
	if search.CustomData == nil {
		search.CustomData = make(map[string]interface{})
	}
	search.CustomData["inputTitle"] = search.InputTitle
}

// Execute the search event, runs the query and sends a search event to
// the analytics. Returns an error if something went wrong.
func (search *SearchEvent) Execute(visit *Visit) (err error) {

	if search.Query == "" { // if the query is empty, randomize one
		var queriesToRandom []string
		if queriesToRandom, err = search.getQueriesToRandomize(visit); err != nil { // Figure out from which queries to randomize
			return
		}
		if search.Query, err = randomQuery(queriesToRandom); err != nil { // Randomize the query from the selected array
			return
		}
	}
	search.Keyword = search.Query
	if search.ActionCause == "" {
		search.ActionCause = defaultSearchCause
	}
	if search.CaseSearch {
		search.handleCaseSearch()
	}
	visit.LastQuery.Q = search.Query
	Info.Printf("Searching for : %s", search.Query)

	// Execute a search and save the response
	if visit.LastResponse, err = visit.SearchClient.Query(*visit.LastQuery); err != nil {
		return
	}

	// in some scenarios (logging of page views), we don't want to send the search event to the analytics
	if !search.IgnoreEvent {
		return search.send(visit)
	}

	Info.Println("Ignoring the search event because of configuration.")
	return
}

func (search *SearchEvent) send(visit *Visit) error {
	if visit.LastResponse == nil {
		return errors.New("LastResponse was nil. Cannot send search event")
	}
	Info.Printf("Sending Search Event with %v results", visit.LastResponse.TotalCount)
	event := ua.NewSearchEvent()

	visit.DecorateEvent(event.ActionEvent)

	event.SearchQueryUID = visit.LastResponse.SearchUID
	event.QueryText = search.Query
	event.AdvancedQuery = visit.LastQuery.AQ
	event.ActionCause = search.ActionCause
	event.NumberOfResults = visit.LastResponse.TotalCount
	event.ResponseTime = visit.LastResponse.Duration

	visit.DecorateCustomMetadata(event.ActionEvent, search.CustomData)

	if visit.LastResponse.TotalCount > 0 {
		if urihash, ok := visit.LastResponse.Results[0].Raw["urihash"].(string); ok {
			event.Results = []ua.ResultHash{
				ua.ResultHash{DocumentURI: visit.LastResponse.Results[0].URI, DocumentURIHash: urihash},
			}
		} else {
			return errors.New("Cannot convert urihash to string in search event")
		}
	}

	// Send a UA search event
	return visit.SendSearchEvent(event)
}

// getQueriesToRandomize Return an array of queries to randomize from.
func (search *SearchEvent) getQueriesToRandomize(visit *Visit) (queriesToRandom []string, err error) {
	if search.GoodQuery { // if we want a good query
		queriesToRandom = visit.Config.GoodQueries
		if search.MatchLanguage { // if the query must match the language
			if _, ok := visit.Config.GoodQueriesInLang[visit.Language]; !ok {
				err = errors.New("No good query detected in " + visit.Language)
				return
			}
			queriesToRandom = visit.Config.GoodQueriesInLang[visit.Language]
		}
	} else { // if we want a bad query
		queriesToRandom = visit.Config.BadQueries
		if search.MatchLanguage { // if the query must match the language
			if _, ok := visit.Config.BadQueriesInLang[visit.Language]; !ok {
				err = errors.New("No bad query detected in " + visit.Language)
				return
			}
			queriesToRandom = visit.Config.BadQueriesInLang[visit.Language]
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
