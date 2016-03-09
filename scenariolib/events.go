// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import (
	"errors"
	"math"
	"math/rand"

	"github.com/k0kubun/pp"
)

// ParseEvent A factory to create the correct event type coming from the JSON parse
// of the scenario definition.
func ParseEvent(e *JSONEvent, c *Config) (Event, error) {
	switch e.Type {

	case "Search":
		event, err := newSearchEvent(e, c)
		if err != nil {
			return nil, err
		}
		return event, nil

	case "Click":
		event, err := newClickEvent(e)
		if err != nil {
			return nil, err
		}
		return event, nil

	case "SearchAndClick":
		event, err := newSearchAndClickEvent(e)
		if err != nil {
			return nil, err
		}
		return event, nil

	case "TabChange":
		event, err := newTabChangeEvent(e)
		if err != nil {
			return nil, err
		}
		return event, nil
	}
	return nil, errors.New("ERR >>> Event type not supported")
}

// Event Generic interface for abstract type Event. All specific event types must
// define the Execute function
type Event interface {
	Execute(v *Visit) error
}

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

// ============== TAB CHANGE EVENT ======================
// ======================================================

// TabChangeEvent represents a tab change event
type TabChangeEvent struct {
	name string
	cq   string
}

func newTabChangeEvent(e *JSONEvent) (*TabChangeEvent, error) {
	name, ok1 := e.Arguments["tabName"].(string)
	cq, ok2 := e.Arguments["tabCQ"].(string)
	if !ok1 || !ok2 {
		return nil, errors.New("ERR >>> Invalid parse of arguments on TabChange Event")
	}

	return &TabChangeEvent{
		name: name,
		cq:   cq,
	}, nil
}

// Execute Sends the tabchange event to the analytics and modify the CQ for the
// following queries in the visit
func (tc *TabChangeEvent) Execute(v *Visit) error {
	pp.Printf("\nLOG >>> Changing tab to %v with CQ : %v", tc.name, tc.cq)

	v.LastQuery.CQ = tc.cq
	v.OriginLevel2 = tc.name
	v.LastQuery.Tab = tc.name

	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}
	v.LastResponse = resp

	pp.Printf("\nLOG >>> Sending Tab Change Event : %v", tc.name)
	err = v.sendInterfaceChangeEvent()
	if err != nil {
		return err
	}
	return nil
}

func (tc *TabChangeEvent) parseTabChangeEvent(e *JSONEvent) error {
	name, ok1 := e.Arguments["tabName"].(string)
	cq, ok2 := e.Arguments["tabCQ"].(string)
	if !ok1 || !ok2 {
		return errors.New("ERR >>> Invalid parse of arguments on TabChange Event")
	}

	tc.name = name
	tc.cq = cq

	return nil
}
