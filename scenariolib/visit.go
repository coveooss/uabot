package scenariolib

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	ua "github.com/coveo/go-coveo/analytics"
	"github.com/coveo/go-coveo/search"
	"github.com/k0kubun/pp"
)

// Visit        The struct visit is used to store one visit to the site.
// SearchClient The http client to send search queries
// UAClient     The http client to send the usage analytics data
// LastQuery    The last query that was searched
// LastResponse The last response that was received
// Username     The name of the user visiting
// OriginLevel1 Where the events originate from
// OriginLevel2 Same as OriginLevel1
// LastTab      The tab the user last visited
type Visit struct {
	SearchClient search.Client
	UAClient     ua.Client
	LastQuery    *search.Query
	LastResponse *search.Response
	Username     string
	OriginLevel1 string
	OriginLevel2 string
	LastTab      string
}

const (
	// JSUIVERSION Change this to the version of JSUI you want to appear to be using.
	JSUIVERSION string = "0.0.0.0;0.0.0.0"
	// TIMEBETWEENACTIONS The time in seconds to wait between the different actions inside a visit
	TIMEBETWEENACTIONS int = 5
)

// NewVisit     Creates a new visit to the search page
// _searchtoken The token used to be able to search
// _uatoken     The token used to send usage analytics events
// _useragent   The user agent the analytics events will see
func NewVisit(_searchtoken string, _uatoken string, _useragent string, c *Config) (*Visit, error) {
	v := Visit{}
	v.Username = fmt.Sprint(c.FirstNames[rand.Intn(len(c.FirstNames))], ".", c.LastNames[rand.Intn(len(c.LastNames))], c.Emails[rand.Intn(len(c.Emails))])
	pp.Printf("\n\nLOG >>> New visit from %v", v.Username)

	// Create the http searchClient
	searchConfig := search.Config{Token: _searchtoken, UserAgent: _useragent, Endpoint: c.SearchEndpoint}
	searchClient, err := search.NewClient(searchConfig)
	if err != nil {
		return nil, err
	}
	v.SearchClient = searchClient

	// Create the http UAClient
	uaConfig := ua.Config{Token: _uatoken, UserAgent: _useragent, IP: c.RandomIPs[rand.Intn(len(c.RandomIPs))], Endpoint: c.AnalyticsEndpoint}
	uaClient, err := ua.NewClient(uaConfig)
	if err != nil {
		return nil, err
	}
	v.UAClient = uaClient

	return &v, nil
}

// ExecuteRandomScenario Method to select randomly a scenario from the possible scenarios and execute it.
// c *Config 	Need the config to have access to the possible random queries and available scenarios
func (v *Visit) ExecuteRandomScenario(c *Config) error {
	scenario, err := c.RandomScenario()
	if err != nil {
		return err
	}

	pp.Printf("\nLOG >>> Executing scenario named : %v", scenario.Name)

	for i := 0; i < len(scenario.Events); i++ {
		event := scenario.Events[i]
		switch event.Type {

		case "Search":
			qText, ok1 := event.Arguments["queryText"].(string)
			goodQuery, ok2 := event.Arguments["goodQuery"].(bool)
			if !ok1 || !ok2 {
				return errors.New("ERR >>> Invalid parse of arguments on Search Event")
			}

			if qText == "" {
				qText, err = c.RandomQuery(goodQuery)
				if err != nil {
					return err
				}
			}

			err := v.executeSearch(qText)
			if err != nil {
				return err
			}

		case "Click":
			offset, ok1 := event.Arguments["offset"].(float64)
			prob, ok2 := event.Arguments["probability"].(float64)
			rank, ok3 := event.Arguments["docNo"].(float64)
			if !ok1 || !ok2 || !ok3 {
				return errors.New("ERR >>> Invalid parse of arguments on Click Event")
			}
			err := v.executeClick(int(rank), int(offset), prob)
			if err != nil {
				return err
			}

		case "SearchAndClick":
			q, ok1 := event.Arguments["queryText"].(string)
			docTitle, ok2 := event.Arguments["docClickTitle"].(string)
			prob, ok3 := event.Arguments["probability"].(float64)
			if !ok1 || !ok2 || !ok3 {
				return errors.New("ERR >>> Invalid parse of arguments on SearchAndClick Event")
			}

			err := v.executeSearch(q)
			if err != nil {
				return err
			}

			if v.LastResponse.TotalCount < 1 {
				return errors.New("ERR >>> Last query returned no results")
			}

			waitBetweenActions()

			if rand.Float64() <= prob {
				rank := v.findDocumentRankByTitle(docTitle)
				if rank >= 0 {
					pp.Printf("\nLOG >>> Sending Click Event => Found Document Rank: %v", rank+1)
					err := v.executeClick(rank, 0, 1)
					if err != nil {
						return err
					}
				} else {
					return errors.New("ERR >>> Could not find the specific document you are looking for")
				}
			} else {
				pp.Printf("\nLOG >>> User chose not to click with a probability of %v %%", int(prob*100))
			}

		case "TabChange":
			name, ok1 := event.Arguments["tabName"].(string)
			cq, ok2 := event.Arguments["tabCQ"].(string)
			if !ok1 || !ok2 {
				return errors.New("ERR >>> Invalid parse of arguments on TabChange Event")
			}
			err := v.executeTabChange(name, cq)
			if err != nil {
				return err
			}
		}
		waitBetweenActions()
	}
	return nil
}

func (v *Visit) executeSearch(q string) error {
	pp.Printf("\nLOG >>> Searching for : %v", q)
	v.LastQuery.Q = q

	// Execute a search and save the response
	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}
	v.LastResponse = resp

	err = v.sendSearchEvent(q)
	if err != nil {
		return err
	}
	return nil
}

func (v *Visit) executeClick(rank int, offset int, prob float64) error {
	if v.LastResponse.TotalCount < 1 {
		pp.Printf("WARN >>> Last query : %v returned no results, cannot click", v.LastQuery.Q)
		return nil
	}

	if prob < 0 || prob > 1 {
		return errors.New("ERR >>> Probability is out of bounds [0..1]")
	}

	if rank == -1 {
		rank = 0
		if v.LastResponse.TotalCount > 1 {
			topL := Min(v.LastQuery.NumberOfResults, v.LastResponse.TotalCount)
			rndRank := int(math.Abs(rand.NormFloat64()*2)) + offset
			rank = Min(rndRank, topL-1)
		}
	}

	if rand.Float64() <= prob {
		if rank > v.LastResponse.TotalCount {
			return fmt.Errorf("\nERR >>> Click Rank out of bounds : %v > %v", rank, v.LastResponse.TotalCount)
		}

		pp.Printf("\nLOG >>> Sending Click Event : Rank -> %v", rank+1)

		err := v.sendClickEvent(rank)
		if err != nil {
			return err
		}
		return nil
	}

	pp.Printf("\nLOG >>> User chose not to click with a probability of : %v %%", int(prob*100))
	return nil
}

func (v *Visit) executeTabChange(name string, cq string) error {
	pp.Printf("\nLOG >>> Changing tab to %v with CQ : %v", name, cq)

	v.LastQuery.CQ = cq
	v.OriginLevel2 = name
	v.LastQuery.Tab = name

	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}
	v.LastResponse = resp

	pp.Printf("\nLOG >>> Sending Tab Change Event : %v", name)
	err = v.sendInterfaceChangeEvent()
	if err != nil {
		return err
	}
	return nil
}

func (v *Visit) sendSearchEvent(q string) error {
	pp.Printf("\nLOG >>> Sending Search Event : %v results", v.LastResponse.TotalCount)
	se, err := ua.NewSearchEvent()
	if err != nil {
		return err
	}

	se.Username = v.Username
	se.SearchQueryUid = v.LastResponse.SearchUID
	se.QueryText = q
	se.AdvancedQuery = v.LastQuery.AQ
	se.ActionCause = "searchboxSubmit"
	se.OriginLevel1 = v.OriginLevel1
	se.OriginLevel2 = v.OriginLevel2
	se.NumberOfResults = v.LastResponse.TotalCount
	se.ResponseTime = v.LastResponse.Duration
	se.CustomData = map[string]interface{}{
		"JSUIVersion": JSUIVERSION,
	}

	if v.LastResponse.TotalCount > 0 {
		if urihash, ok := v.LastResponse.Results[0].Raw["sysurihash"].(string); ok {
			se.Results = []ua.ResultHash{
				ua.ResultHash{DocumentUri: v.LastResponse.Results[0].URI, DocumentUriHash: urihash},
			}
		} else {
			return errors.New("ERR >>> Cannot convert sysurihash to string in search event")
		}
	}

	// Send a UA search event
	error := v.UAClient.SendSearchEvent(se)
	if error != nil {
		return err
	}
	return nil
}

func (v *Visit) sendClickEvent(rank int) error {
	event, err := ua.NewClickEvent()
	if err != nil {
		return err
	}

	event.DocumentUri = v.LastResponse.Results[rank].URI
	event.SearchQueryUid = v.LastResponse.SearchUID
	event.DocumentPosition = rank + 1
	event.ActionCause = "documentOpen"
	event.DocumentTitle = v.LastResponse.Results[rank].Title
	event.QueryPipeline = v.LastResponse.Pipeline
	event.DocumentUrl = v.LastResponse.Results[rank].ClickUri
	event.Username = v.Username
	event.OriginLevel1 = v.OriginLevel1
	event.OriginLevel2 = v.OriginLevel2
	if urihash, ok := v.LastResponse.Results[rank].Raw["sysurihash"].(string); ok {
		event.DocumentUriHash = urihash
	} else {
		return errors.New("ERR >>> Cannot convert sysurihash to string")
	}
	if collection, ok := v.LastResponse.Results[rank].Raw["syscollection"].(string); ok {
		event.CollectionName = collection
	} else {
		return errors.New("ERR >>> Cannot convert syscollection to string")
	}
	if source, ok := v.LastResponse.Results[rank].Raw["syssource"].(string); ok {
		event.SourceName = source
	} else {
		return errors.New("ERR >>> Cannot convert syssource to string")
	}

	err = v.UAClient.SendClickEvent(event)
	if err != nil {
		return err
	}
	return nil
}

func (v *Visit) sendInterfaceChangeEvent() error {
	ice, err := ua.NewSearchEvent()
	if err != nil {
		return err
	}

	ice.Username = v.Username
	ice.SearchQueryUid = v.LastResponse.SearchUID
	ice.QueryText = v.LastQuery.Q
	ice.AdvancedQuery = v.LastQuery.AQ
	ice.ActionCause = "interfaceChange"
	ice.OriginLevel1 = v.OriginLevel1
	ice.OriginLevel2 = v.OriginLevel2
	ice.NumberOfResults = v.LastResponse.TotalCount
	ice.ResponseTime = v.LastResponse.Duration
	ice.CustomData = map[string]interface{}{
		"interfaceChangeTo": v.OriginLevel2,
		"JSUIVersion":       JSUIVERSION,
	}

	if v.LastResponse.TotalCount > 0 {
		if urihash, ok := v.LastResponse.Results[0].Raw["sysurihash"].(string); ok {
			ice.Results = []ua.ResultHash{
				ua.ResultHash{DocumentUri: v.LastResponse.Results[0].URI, DocumentUriHash: urihash},
			}
		} else {
			return errors.New("ERR >>> Cannot convert sysurihash to string in interfaceChange event")
		}
	}

	err = v.UAClient.SendSearchEvent(ice)
	if err != nil {
		return err
	}
	return nil
}

func (v *Visit) findDocumentRankByTitle(toFind string) int {
	for i := 0; i < len(v.LastResponse.Results); i++ {
		if strings.Contains(strings.ToLower(v.LastResponse.Results[i].Title), strings.ToLower(toFind)) {
			return i
		}
	}
	return -1
}

// Wait a random number of seconds between user actions
func waitBetweenActions() {
	time.Sleep(time.Duration(rand.Intn(TIMEBETWEENACTIONS)+3) * time.Second)
}

// Min Function to return the minimal value between two integers, because Go "forgot"
// to code it...
func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// SetupNTO Function to instanciate with specific values for NTO demo queries
func (v *Visit) SetupNTO() {
	gbs := []*search.GroupByRequest{}
	q := &search.Query{
		Q:               "",
		CQ:              "",
		AQ:              "NOT @objecttype==(User,Case,CollaborationGroup) AND NOT @sysfiletype==(Folder, YouTubePlaylist, YouTubePlaylistItem)",
		NumberOfResults: 20,
		FirstResult:     0,
		Tab:             "All",
		GroupByRequests: gbs,
	}

	v.LastQuery = q

	v.OriginLevel1 = "communityCoveo"
	v.OriginLevel2 = "ALL"
}

// SetupGeneral Function to instanciate with non-specific values
func (v *Visit) SetupGeneral() {
	gbs := []*search.GroupByRequest{}
	q := &search.Query{
		Q:               "",
		CQ:              "",
		AQ:              "",
		NumberOfResults: 20,
		FirstResult:     0,
		Tab:             "All",
		GroupByRequests: gbs,
	}

	v.LastQuery = q

	v.OriginLevel1 = "ALL"
	v.OriginLevel2 = "ALL"
}
