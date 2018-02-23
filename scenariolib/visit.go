package scenariolib

import (
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	ua "github.com/coveo/go-coveo/analytics"
	"github.com/coveo/go-coveo/search"
	"github.com/coveo/uabot/defaults"
)

// Visit        The struct visit is used to store one visit to the site.
// SearchClient The http client to send search queries
// UAClient     The http client to send the usage analytics data
// LastQuery    The last query that was searched
// LastResponse The last response that was received
// Username     The name of the user visiting
// OriginLevel1 Page/Hub where the events originate from
// OriginLevel2 Tab where the events originate from
// OriginLevel3 The HTTP identifier of the page from which any type of event originates
// Referrer     Same as OriginLevel3
// LastTab      The tab the user last visited
type Visit struct {
	SearchClient       search.Client
	UAClient           ua.Client
	LastQuery          *search.Query
	LastResponse       *search.Response
	Username           string
	OriginLevel1       string
	OriginLevel2       string
	OriginLevel3       string
	Referrer           string
	LastTab            string
	Config             *Config
	IP                 string
	Anonymous          bool
	Language           string
	WaitBetweenActions bool
}

const (
	// JSUIVERSION Change this to the version of JSUI you want to appear to be using.
	JSUIVERSION string = "0.0.0.0;0.0.0.0"
	// DEFAULTTIMEBETWEENACTIONS The time in seconds to wait between the different actions inside a visit
	DEFAULTTIMEBETWEENACTIONS int = 5
	// ORIGINALL The origin level of All
	ORIGINALL string = "ALL"
)

// NewVisit     Creates a new visit to the search page
// _searchtoken The token used to be able to search
// _uatoken     The token used to send usage analytics events
// _useragent   The user agent the analytics events will see
func NewVisit(_searchtoken string, _uatoken string, _useragent string, language string, c *Config) (*Visit, error) {

	InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	v := Visit{}
	v.Config = c

	v.WaitBetweenActions = !c.DontWaitBetweenVisits
	v.Anonymous = false

	if c.AnonymousThreshold > 0 {
		if rand.Float64() <= c.AnonymousThreshold {
			v.Anonymous = true
			Info.Printf("Anonymous visit")
		}
	}

	if !v.Anonymous {
		v.Username = buildUserEmail(c)
		Info.Printf("New visit from %s", v.Username)
	}

	if language != "" {
		v.Language = language
	} else {
		if len(v.Config.RandomData.Languages) > 0 {
			v.Language = randomStringArray(v.Config.RandomData.Languages)
		} else {
			v.Language = "en"
		}
	}
	Info.Printf("Language of visit : %s", v.Language)

	// Create the http searchClient
	searchConfig := search.Config{Token: _searchtoken, UserAgent: _useragent, Endpoint: c.SearchEndpoint}
	searchClient, err := search.NewClient(searchConfig)
	if err != nil {
		return nil, err
	}
	v.SearchClient = searchClient

	// Create the http UAClient
	v.IP = randomStringArray(c.RandomData.RandomIPs)
	uaConfig := ua.Config{Token: _uatoken, UserAgent: _useragent, IP: v.IP, Endpoint: c.AnalyticsEndpoint}
	uaClient := ua.NewClient(uaConfig)
	v.UAClient = uaClient

	return &v, nil
}

func buildUserEmail(c *Config) string {
	firstName := randomStringArray(c.RandomData.FirstNames)
	lastName := randomStringArray(c.RandomData.LastNames)
	email := randomStringArray(c.RandomData.Emails)
	return fmt.Sprint(firstName, ".", lastName, email)
}

// randomStringArray Returns a random string value from a []string
func randomStringArray(array []string) (value string) {
	value = array[rand.Intn(len(array))]
	return
}

// ExecuteScenario Execute a specific scenario, send the config for all the
// potential random we need to do.
func (v *Visit) ExecuteScenario(scenario Scenario, c *Config) error {
	Info.Printf("Executing scenario named : %s", scenario.Name)
	for i := 0; i < len(scenario.Events); i++ {
		jsonEvent := scenario.Events[i]
		event, err := ParseEvent(&jsonEvent, c)
		if err != nil {
			return err
		}
		err = event.Execute(v)
		if err != nil {
			return err
		}
		if v.WaitBetweenActions {
			if c.TimeBetweenActions > 0 {
				WaitBetweenActions(c.TimeBetweenActions, c.IsWaitConstant)
			} else {
				WaitBetweenActions(DEFAULTTIMEBETWEENACTIONS, c.IsWaitConstant)
			}

		}
	}
	return nil
}

// SendSearchEvent to the UAClient.
func (v *Visit) SendSearchEvent(event *ua.SearchEvent) (err error) {
	err = v.UAClient.SendSearchEvent(event)
	return
}

// SendViewEvent to the UAClient
func (v *Visit) SendViewEvent(event *ua.ViewEvent) (err error) {
	err = v.UAClient.SendViewEvent(event)
	return
}

func (v *Visit) sendCustomEvent(actionCause, actionType string, customData map[string]interface{}) error {
	Info.Printf("Sending CustomEvent cause: %s ||| type: %s", actionCause, actionType)
	event := ua.NewCustomEvent()

	v.DecorateEvent(event.ActionEvent)

	event.EventType = actionType
	event.EventValue = actionCause

	v.DecorateCustomMetadata(event.ActionEvent, customData)

	if v.Config.AllowEntitlements {
		// Custom shit for besttech, I don't like it
		event.CustomData["entitlement"] = generateEntitlementBesttech(v.Anonymous)
	}

	// Send a UA search event
	err := v.UAClient.SendCustomEvent(event)
	return err
}

func (v *Visit) sendClickEvent(rank int, quickview bool, customData map[string]interface{}) error {
	if v.LastResponse == nil {
		return errors.New("LastResponse was nil cannot send click event")
	}
	Info.Printf("Sending ClickEvent rank=%d (quickview %v)", rank+1, quickview)
	event := ua.NewClickEvent()

	v.DecorateEvent(event.ActionEvent)

	event.SearchQueryUID = v.LastResponse.SearchUID
	event.DocumentURI = v.LastResponse.Results[rank].URI
	event.DocumentTitle = v.LastResponse.Results[rank].Title
	event.QueryPipeline = v.LastResponse.Pipeline
	event.DocumentURL = v.LastResponse.Results[rank].ClickURI
	event.DocumentPosition = rank + 1 //Document Position is 1 based in UA

	if quickview {
		event.ActionCause = "documentQuickview"
	} else {
		event.ActionCause = "documentOpen"
	}

	if urihash, ok := v.LastResponse.Results[rank].Raw["urihash"].(string); ok {
		event.DocumentURIHash = urihash
	} else {
		return errors.New("Cannot convert urihash to string")
	}
	if collection, ok := v.LastResponse.Results[rank].Raw["collection"].(string); ok {
		event.CollectionName = collection
	} else {
		// TODO: handle indexless option here
		event.CollectionName = "default"
		Warning.Println("Cannot convert collection to string, sending \"default\"")
		//return errors.New("Cannot convert collection to string")
	}
	if source, ok := v.LastResponse.Results[rank].Raw["source"].(string); ok {
		event.SourceName = source
	} else {
		// TODO: handle indexless option here
		event.SourceName = "default"
		Warning.Println("Cannot convert source to string, sending \"default\"")
		// return errors.New("Cannot convert source to string")
	}

	v.DecorateCustomMetadata(event.ActionEvent, customData)

	if v.Config.AllowEntitlements {
		// Custom shit for besttech, I don't like it
		event.CustomData["entitlement"] = generateEntitlementBesttech(v.Anonymous)
	}

	event.CustomData["author"] = generateRandomAuthor(event.DocumentTitle)

	err := v.UAClient.SendClickEvent(event)
	if err != nil {
		return err
	}
	return nil
}

func (v *Visit) sendInterfaceChangeEvent(actionCause, actionType string, customData map[string]interface{}) error {
	if v.LastResponse == nil {
		return errors.New("LastResponse was nil cannot send InterfaceChange event")
	}
	event := ua.NewSearchEvent()

	// Add all the metadata on the event that is common across all events.
	v.DecorateEvent(event.ActionEvent)

	event.SearchQueryUID = v.LastResponse.SearchUID
	event.QueryText = v.LastQuery.Q
	event.AdvancedQuery = v.LastQuery.AQ
	event.ActionCause = actionCause
	event.NumberOfResults = v.LastResponse.TotalCount
	event.ResponseTime = v.LastResponse.Duration

	if v.LastResponse.TotalCount > 0 {
		if urihash, ok := v.LastResponse.Results[0].Raw["urihash"].(string); ok {
			event.Results = []ua.ResultHash{
				ua.ResultHash{DocumentURI: v.LastResponse.Results[0].URI, DocumentURIHash: urihash},
			}
		} else {
			return errors.New("Cannot convert urihash to string in interfaceChange event")
		}
	}

	v.DecorateCustomMetadata(event.ActionEvent, customData)

	if v.Config.AllowEntitlements {
		// Custom shit for besttech, I don't like it
		event.CustomData["entitlement"] = generateEntitlementBesttech(v.Anonymous)
	}

	err := v.UAClient.SendSearchEvent(event)
	if err != nil {
		return err
	}
	return nil
}

// DecorateEvent is used to assign all the common data to send with all analytics events
func (v *Visit) DecorateEvent(evt *ua.ActionEvent) {
	evt.Username = v.Username
	evt.Anonymous = v.Anonymous
	evt.Language = v.Language

	evt.OriginLevel1 = v.OriginLevel1
	evt.OriginLevel2 = v.OriginLevel2
	evt.OriginLevel3 = v.OriginLevel3

	// Only send AB tests metadata if there is an acutal AB test involved in the last Response.
	if v.LastResponse.SplitTestRun != "" {
		evt.SplitTestRunName = v.LastResponse.SplitTestRun
		evt.SplitTestRunVersion = v.LastResponse.Pipeline
	}
}

// DecorateCustomMetadata is used to handle all the customMetadata for the events that allow it.
func (v *Visit) DecorateCustomMetadata(evt *ua.ActionEvent, customData map[string]interface{}) {

	evt.CustomData = map[string]interface{}{
		"JSUIVersion": JSUIVERSION,
		"ipaddress":   v.IP,
	}

	// Send all the possible random custom data that can be added from the config
	// scenario file.
	for _, elem := range v.Config.RandomCustomData {
		evt.CustomData[elem.APIName] = randomStringArray(elem.Values)
	}

	// Override possible values of customData with the specific customData sent
	for k, v := range customData {
		evt.CustomData[k] = v
	}
}

// FindDocumentRankByTitle Looks through the last response to a query to find a document
// rank by his title
func (v *Visit) FindDocumentRankByTitle(toFind string) int {
	if v.LastResponse == nil {
		return -1
	}
	for i := 0; i < len(v.LastResponse.Results); i++ {
		if strings.Contains(strings.ToLower(v.LastResponse.Results[i].Title), strings.ToLower(toFind)) {
			return i
		}
	}
	return -1
}

// WaitBetweenActions Wait a random or constant number of seconds between user actions
func WaitBetweenActions(timeToWait int, isConstant bool) {
	timeToWait = timeToWait * 1000
	if !isConstant {
		timeToWait = rand.Intn(timeToWait)
	}
	timeToWait = Max(timeToWait, 500)

	time.Sleep(time.Duration(timeToWait) * time.Millisecond)
}

// Min Function to return the minimal value between two integers, because Go "forgot"
// to code it...
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max function returning the maximum between two values of type int
func Max(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func generateRandomAuthor(title string) string {
	return defaults.AUTHORNAMES[(int)(math.Mod((float64)(hash(title)), (float64)(len(defaults.AUTHORNAMES))))]
}

func generateEntitlementBesttech(isAnonymous bool) string {
	if isAnonymous {
		return "Anonymous"
	}
	if rand.Float64() <= 0.1 {
		return "Premier"
	}
	return "Basic"
}

// SetupNTO Function to instanciate with specific values for NTO demo queries
func (v *Visit) SetupNTO() {
	gbs := []*search.GroupByRequest{}
	q := &search.Query{
		Q:               "",
		CQ:              "",
		AQ:              "NOT @objecttype==(User,Case,CollaborationGroup) AND NOT @filetype==(Folder, YouTubePlaylist, YouTubePlaylistItem)",
		NumberOfResults: 20,
		FirstResult:     0,
		Tab:             "All",
		GroupByRequests: gbs,
	}

	if v.Config.PartialMatch {
		q.PartialMatch = v.Config.PartialMatch
		q.PartialMatchKeywords = v.Config.PartialMatchKeywords
		q.PartialMatchThreshold = v.Config.PartialMatchThreshold
	}

	if v.Config.Pipeline != "" {
		q.Pipeline = v.Config.Pipeline
	}

	v.LastQuery = q

	v.OriginLevel1 = "Community"
	v.OriginLevel2 = ORIGINALL
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

	if v.Config.PartialMatch {
		q.PartialMatch = v.Config.PartialMatch
		q.PartialMatchKeywords = v.Config.PartialMatchKeywords
		q.PartialMatchThreshold = v.Config.PartialMatchThreshold
	}

	if v.Config.Pipeline != "" {
		q.Pipeline = v.Config.Pipeline
	}

	v.LastQuery = q

	v.OriginLevel1 = v.Config.DefaultOriginLevel1
	if v.Config.DefaultOriginLevel2 != "" {
		v.OriginLevel2 = v.Config.DefaultOriginLevel2
	}
	if v.Config.DefaultOriginLevel3 != "" {
		v.OriginLevel3 = v.Config.DefaultOriginLevel3
	}
}
