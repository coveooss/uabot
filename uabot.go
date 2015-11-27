package main

import (
	"errors"
	"fmt"
	ua "github.com/coveo/go-coveo/analytics"
	"github.com/coveo/go-coveo/search"
	"github.com/k0kubun/pp"
	"math/rand"
	"math"
	"time"
	"strings"
	"flag"
	"encoding/json"
	"net/http"
	"os"
)

const (
	DEBUG              int    = 0
	JSUIVERSION        string = "0.0.0.0;0.0.0.0"
	USERAGENT          string = "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36"
	//NUMBEROFCASES      int    = 100
	TIMEBETWEENVISITS  int    = 120 // Between 0 and X Seconds
	TIMEBETWEENACTIONS int    = 10 // Between 0 and X Seconds
)

/*var RANDOMEMAILS []string = []string{
	"@gmail.com", "@hotmail.com", "@apple.com", "@yahoo.com", "@facebook.com", "@hexzone.com", "@strongit.com", "@mec.com",
	"@geoflex.com",
}

var RANDOMFIRSTNAMES []string = []string{
	"erin", "paul", "beverley", "pedro", "clayton", "lydia", "regina", "sue", "marjorie", "april", "victoria", "vera",
	"shannon", "minnie", "reginald", "brandie",	"christian", "wallace", "avery", "dawn",
}

var RANDOMLASTNAMES []string = []string{
	"lawson", "torres", "grant", "ray", "young", "caldwell", "morris", "craig",	"lewis", "brown", "rhodes", "james",
	"wagner", "richards", "allen", "berry",	"boyd", "price", "price", "rivera",
}

var RDQUERIES []string = []string{
	"tent", "bike helmet", "hiking", "race of hope", "hiking equipment", "backpack", "flat tires",
	"runner pack", "hiking boots", "user manual", "bag", "camping gear", "trails outfitter",
	"how to repair my tent", "rivendale 800", "kampa tent", "mec tent", "hiking tips",
	"climbing shoe", "camping shelter", "darkstar", "running equipment", "tire pressure", "repair flat tire",
	"tent assembly", "tour backpack", "best family tents", "yeti tundra", "portable usb charger",
	"waterproof jacket", "saddles", "airbeam tent", "checkout",
}

var RDBADQUERIES []string = []string {
	"rivendale 801", "kampa explorer 8", "stuck zip", "my zipper is broken", "tents", "camping", "travel", "bike",
	"hike", "northern trail", "northern trail outfitters", "northern trails", "travel tips", "hiking and tent",
	"travel tip",
}

var RANDOMIPS []string = []string{
	"66.46.18.120", "74.125.226.120", "66.46.18.1", "192.40.239.233", // Canada
	"198.169.156.67", "160.72.0.1", "155.15.0.45", "162.248.127.25", // Canada
	"52.24.0.108", "159.28.0.98", "205.214.160.167", "216.252.192.109", // US
	"72.9.32.109", "198.199.154.209", "209.137.0.105", "216.249.112.8", // US
}*/

var SearchToken string = ""
var UAToken     string = ""

var ScenariosDef ScenariosDefinition = ScenariosDefinition{}

type UseCase struct {
	Cs           search.Client
	Cua          ua.Client
	LastQuery    search.Query
	LastResponse *search.Response
	Username     string
	OriginLevel1 string
	OriginLevel2 string
	Debug        int
}

type ScenariosDefinition struct {
	OrgName     string      `json:"orgName"`
	Emails      []string    `json:"emailSuffixes"`
	FirstNames  []string    `json:"firstNames"`
	LastNames   []string    `json:"lastNames"`
	GoodQueries []string    `json:"randomGoodQueries"`
	BadQueries  []string    `json:"randomBadQueries"`
	RandomIPs   []string    `json:"randomIPs"`
	Scenarios   []*Scenario `json:"scenarios"`
}

type Scenario struct {
	Name   string  `json:"name"`
	Weight int     `json:"weight"`
	Events []Event `json:"events"`
}

type Event struct {
	Type      string                 `json:"type"`
	Arguments map[string]interface{} `json:"arguments"`
}

func newUseCase() (*UseCase, error) {
	// Create the Search client.
	conf_s := search.Config{Token: SearchToken, UserAgent: USERAGENT}
	cs, err := search.NewClient(conf_s)
	if err != nil {
		return nil, err
	}

	// Create the UA client.
	conf_ua := ua.Config{Token: UAToken, UserAgent: USERAGENT, IP: ScenariosDef.RandomIPs[rand.Intn(len(ScenariosDef.RandomIPs))], }
	cua, err := ua.NewClient(conf_ua)
	if err != nil {
		return nil, err
	}

	return &UseCase{
		Cs:  cs,
		Cua: cua,
	}, nil
}

func sendSearchEvent(useCase *UseCase) error {
	SEvent, err := ua.NewSearchEvent()
	if err != nil { return err }

	SEvent.Username = useCase.Username
	SEvent.SearchQueryUid = useCase.LastResponse.SearchUID
	SEvent.QueryText = useCase.LastQuery.Q
	SEvent.AdvancedQuery = useCase.LastQuery.AQ
	SEvent.ActionCause = "searchboxSubmit"
	SEvent.OriginLevel1 = useCase.OriginLevel1
	SEvent.OriginLevel2 = useCase.OriginLevel2
	SEvent.NumberOfResults = useCase.LastResponse.TotalCount
	SEvent.ResponseTime = useCase.LastResponse.Duration
	SEvent.CustomData = map[string]interface{}{
		"JSUIVersion" : JSUIVERSION,
	}

	if useCase.LastResponse.TotalCount > 0 {
		if urihash, ok := useCase.LastResponse.Results[0].Raw["sysurihash"].(string); ok {
			SEvent.Results = []ua.ResultHash{
				ua.ResultHash{DocumentUri:useCase.LastResponse.Results[0].URI, DocumentUriHash: urihash},
			}
		} else { return errors.New("Cannot convert sysurihash to string") }
	}

	error := useCase.Cua.SendSearchEvent(SEvent)
	if error != nil { return err }

	return nil
}

func sendInterfaceChangeEvent(useCase *UseCase) error {
	ICEvent, err := ua.NewSearchEvent()
	if err != nil { return err }

	ICEvent.Username = useCase.Username
	ICEvent.SearchQueryUid = useCase.LastResponse.SearchUID
	ICEvent.QueryText = useCase.LastQuery.Q
	ICEvent.AdvancedQuery = useCase.LastQuery.AQ
	ICEvent.ActionCause = "interfaceChange"
	ICEvent.OriginLevel1 = useCase.OriginLevel1
	ICEvent.OriginLevel2 = useCase.OriginLevel2
	ICEvent.NumberOfResults = useCase.LastResponse.TotalCount
	ICEvent.ResponseTime = useCase.LastResponse.Duration
	ICEvent.CustomData = map[string]interface{}{
		"interfaceChangeTo" : useCase.OriginLevel2,
		"JSUIVersion"       : JSUIVERSION,
	}

	if useCase.LastResponse.TotalCount > 0 {
		if urihash, ok := useCase.LastResponse.Results[0].Raw["sysurihash"].(string); ok {
			ICEvent.Results = []ua.ResultHash{
				ua.ResultHash{DocumentUri:useCase.LastResponse.Results[0].URI, DocumentUriHash: urihash},
			}
		} else { return errors.New("Cannot convert sysurihash to string") }
	}

	error := useCase.Cua.SendSearchEvent(ICEvent)
	if error != nil { return err }

	return nil
}

func sendClickDocument(useCase *UseCase, docNo int) error {

	resp := useCase.LastResponse

	if docNo > resp.TotalCount {
		return errors.New("Cannot click on a document that doesn't exist")
	}

	CEvent, err := ua.NewClickEvent()
	if err != nil {
		return err
	}

	CEvent.DocumentUri = resp.Results[docNo].URI
	if urihash, ok := resp.Results[docNo].Raw["sysurihash"].(string); ok {
		CEvent.DocumentUriHash = urihash
	} else { return errors.New("Cannot convert sysurihash to string") }

	CEvent.SearchQueryUid = resp.SearchUID
	if collection, ok := resp.Results[docNo].Raw["syscollection"].(string); ok {
		CEvent.CollectionName = collection
	} else { return errors.New("Cannot convert syscollection to string") }

	if source, ok := resp.Results[docNo].Raw["syssource"].(string); ok {
		CEvent.SourceName = source
	} else { return errors.New("Cannot convert syssource to string") }

	CEvent.DocumentPosition = docNo
	CEvent.ActionCause = "documentOpen"
	CEvent.DocumentTitle = resp.Results[docNo].Title
	CEvent.QueryPipeline = resp.Pipeline
	CEvent.DocumentUrl = resp.Results[docNo].ClickUri
	CEvent.Username = useCase.Username
	CEvent.OriginLevel1 = useCase.OriginLevel1
	CEvent.OriginLevel2 = useCase.OriginLevel2

	error := useCase.Cua.SendClickEvent(CEvent)
	if error != nil {
		return error
	}

	return nil
}

func ntoSetupFirstQuery(useCase *UseCase) error {

	// ==================================================
	// Setup the first query to init the search interface
	// ==================================================
	gb1 := search.GroupByRequest{Field: "@syssource", MaximumNumberOfValues: 6, SortCriteria: "occurences", InjectionDepth: 1000}
	gb2 := search.GroupByRequest{Field: "@coveochatterfeedtopics", MaximumNumberOfValues: 6, SortCriteria: "occurences", InjectionDepth: 1000}
	gb3 := search.GroupByRequest{Field: "@sysyear", MaximumNumberOfValues: 6, SortCriteria: "occurences", InjectionDepth: 1000}
	gbs := []*search.GroupByRequest{&gb1, &gb2, &gb3}

	q := search.Query{
		Q:               "",
		CQ:              "",
		AQ:              "NOT @objecttype==(User,Case,CollaborationGroup) AND NOT @sysfiletype==(Folder, YouTubePlaylist, YouTubePlaylistItem)",
		NumberOfResults: 20,
		FirstResult:     0,
		Tab:             "All",
		GroupByRequests: gbs,
	}

	useCase.LastQuery = q

	useCase.OriginLevel1 = "communityCoveo"
	useCase.OriginLevel2 = "ALL"

	return nil
}

// Execute a query and send a search event to UA
//
// queryText  Leave empty if you want a random query
// goodQuery  If true will search in the pool of good queries
func NewSearchUseCase(useCase *UseCase, queryText string, goodQuery bool) error {

	if queryText == "" {
		if goodQuery { queryText = ScenariosDef.GoodQueries[Min(int(math.Abs(rand.NormFloat64() * 8)), len(ScenariosDef.GoodQueries))] 
		} else { queryText = ScenariosDef.BadQueries[rand.Intn(len(ScenariosDef.BadQueries))] }
	}

	useCase.LastQuery.Q = queryText

	if useCase.Debug == 1 { pp.Printf(">> Searching for : %v \n", useCase.LastQuery.Q) }

	qResponse, err := useCase.Cs.Query(useCase.LastQuery)
	if err != nil { return err }

	useCase.LastResponse = qResponse

	if useCase.Debug == 1 { pp.Printf(">> Sending Search Event >>> : %v results\n", useCase.LastResponse.TotalCount) }

	err = sendSearchEvent(useCase)
	if err != nil { return err }

	return nil
}

// Generate a click event to UA
//
// docNo        The rank of the document to click, -1 for a random click
// offset       Offset of the random click rank
// probability  Probability of the click happening (0.30 means 30 chance to click)
func NewClickUseCase(useCase *UseCase, docNo int, offset int, probability float64) error {
	if useCase.LastResponse.TotalCount > 0 { 
		if probability < 0 || probability > 1 { return errors.New("probability is lower than 0 or greater than 1") }

		if docNo == -1 {
			docNo = 0
			if useCase.LastResponse.TotalCount > 1 {
				topLimit := Min(useCase.LastQuery.NumberOfResults, useCase.LastResponse.TotalCount)
				randomRank := int(math.Abs(rand.NormFloat64() * 2)) + offset
				docNo = Min(randomRank, topLimit-1)
			}
		}

		randProb := rand.Float64()

		if randProb <= probability {
			if useCase.Debug == 1 { pp.Printf(">> Sending Click Event >>> Document Rank: %v\n", docNo) }
			err := sendClickDocument(useCase, docNo)
			if err != nil { return err }
		} else if useCase.Debug == 1 {
			pp.Println(">> User chose to not click")
		}
	}
	return nil
}

// Execute a query, send a search event to UA and find a specific document to click on to send click event to UA
//
// queryText     The text of the query to execute
// docClickTitle The title of the document to click on
// probability   Probability of the click to happen 0-1
func NewSearchAndClickUseCase(useCase *UseCase, queryText string, docClickTitle string, probability float64) error {
	if probability < 0 || probability > 1 { return errors.New("probability is lower than 0 or greater than 1") }

	useCase.LastQuery.Q = queryText

	if useCase.Debug == 1 { pp.Printf(">> Searching for : %v \n", useCase.LastQuery.Q) }

	qResponse, err := useCase.Cs.Query(useCase.LastQuery)
	if err != nil { return err }

	useCase.LastResponse = qResponse

	if useCase.Debug == 1 { pp.Printf(">> Sending Search Event >>> : %v results\n", useCase.LastResponse.TotalCount) }

	err = sendSearchEvent(useCase)
	if err != nil { return err }

	WaitBetweenActions()

	if useCase.LastResponse.TotalCount < 1 { return errors.New("Last query returned no results") }

	randProb := rand.Float64()

	if randProb <= probability {
		docNo := FindDocumentRankByTitle(docClickTitle, useCase.LastResponse.Results)
		if docNo >= 0 {
			if useCase.Debug == 1 { pp.Printf(">> Sending Click Event >>> Found Document Rank: %v\n", docNo) }
			
			err := sendClickDocument(useCase, docNo)
			if err != nil { return err }
		} else {
			//TO-DO Better error handling
			pp.Printf("!! Could not find document titled >>> %v in the results of the query\n", docClickTitle)
			pp.Printf("!! Query >> %v\n", useCase.LastQuery.Q)
		}
	} else if useCase.Debug == 1 {
		pp.Println(">> User chose not to click")
	}

	return nil
}

// Execute a tab change event, change the CQ, execute a query and send a search event
//
// tabName The name of the tab to change to
// tabCQ   The CQ associated with the tab
func NewTabChangeUseCase(useCase *UseCase, tabName string, tabCQ string) error {
	if useCase.Debug == 1 { pp.Printf(">> Changing tab to : %v >>> CQ : %v\n", tabName, tabCQ) }

	useCase.LastQuery.CQ = tabCQ
	useCase.OriginLevel2 = tabName

	qResponse, err := useCase.Cs.Query(useCase.LastQuery)
	if err != nil { return err }

	useCase.LastResponse = qResponse

	if useCase.Debug == 1 { pp.Printf(">> Sending Tab Change Event >>> : %v \n", tabName) }

	err = sendInterfaceChangeEvent(useCase)
	if err != nil { return err }

	return nil
}

func ExecuteScenario(scenario *Scenario, useCase *UseCase) error {
	if useCase.Debug == 1 { pp.Printf("\n>> Executing Scenario >>> %v\n", scenario.Name) }
	for i := 0; i < len(scenario.Events); i++ {
		event := scenario.Events[i]

		//if useCase.Debug == 1 { pp.Println(scenario.Events[i]) }
		switch event.Type {

		case "Search" :
			queryText, ok1 := event.Arguments["queryText"].(string)
			goodQuery, ok2 := event.Arguments["goodQuery"].(bool)
			if ok1 && ok2 {
				err := NewSearchUseCase(useCase, queryText, goodQuery)
				if err != nil { return err }
			} else { return errors.New("!! Invalid parse for event Search >>> Cannot read arguments") }

		case "Click" :
			offset, ok1 := event.Arguments["offset"].(float64)
			probability, ok2 := event.Arguments["probability"].(float64)
			docNo, ok3 := event.Arguments["docNo"].(float64)
			if ok1 && ok2 && ok3 {
				err := NewClickUseCase(useCase, int(docNo), int(offset), probability)
				if err != nil { return err }
			} else { return errors.New("!! Invalid parse for event Click >>> Cannot read arguments") }

		case "SearchAndClick" :
			queryText, ok1 := event.Arguments["queryText"].(string)
			docClickTitle, ok2 := event.Arguments["docClickTitle"].(string)
			probability, ok3 := event.Arguments["probability"].(float64)
			if ok1 && ok2 && ok3 {
				err := NewSearchAndClickUseCase(useCase, queryText, docClickTitle, probability)
				if err != nil { return err }
			} else { return errors.New("!! Invalid parse for event SearchAndClick >>> Cannot read arguments") }
		case "TabChange" :
			tabName, ok1 := event.Arguments["tabName"].(string)
			tabCQ, ok2 := event.Arguments["tabCQ"].(string)
			if ok1 && ok2 {
				err := NewTabChangeUseCase(useCase, tabName, tabCQ)
				if err != nil { return err }
			} else { return errors.New("!! Invalid parse for event TabChange >>> Cannot read arguments") }
		}
		WaitBetweenActions()
	}
	return nil
}

func InitUseCase(debug int) (*UseCase, error) {

	useCase, err := newUseCase()
	if err != nil { return nil, err }

	useCase.Debug = debug

	// Randomize a username
	useCase.Username = fmt.Sprint(ScenariosDef.FirstNames[rand.Intn(len(ScenariosDef.FirstNames))], ".", ScenariosDef.LastNames[rand.Intn(len(ScenariosDef.LastNames))], ScenariosDef.Emails[rand.Intn(len(ScenariosDef.Emails))])

	if useCase.Debug == 1 {
		pp.Printf("\n===============================================")
		pp.Printf("\n>> New Visit >>> Username : %v\n", useCase.Username)
	}

	// Starting query
	err = ntoSetupFirstQuery(useCase)
	if err != nil { return nil, err }

	return useCase, nil
}

func ParseScenariosFile(url string) (map[int]*Scenario, error) {

	resp, err := http.Get(url)
	if err != nil { return nil, err }
	defer resp.Body.Close()

	scenariosDef := &ScenariosDefinition{}
	err = json.NewDecoder(resp.Body).Decode(&scenariosDef)
	if err != nil { return nil, err }

	ScenariosDef = *scenariosDef

	var scenarioMap map[int]*Scenario = map[int]*Scenario{}

	totalWeight := 0
	iter := 0
	for i := 0; i < len(scenariosDef.Scenarios); i++ {
		weight := scenariosDef.Scenarios[i].Weight
		totalWeight += weight
		for j := 0; j < weight; j++ {
			scenarioMap[iter] = scenariosDef.Scenarios[i]
			iter++
		}
	}

	return scenarioMap, nil
}

func main() {

	debug := flag.Int("debug", 0, "DEBUG MODE")
	flag.Parse()

	rand.Seed(int64(time.Now().Unix()))

	sToken := os.Getenv("SEARCHTOKEN")
	uaToken := os.Getenv("UATOKEN")
	if sToken == "" || uaToken == "" { pp.Fatal("No search token or UA token") }

	SearchToken = sToken
	UAToken = uaToken

	scenarioUrl := os.Getenv("SCENARIOSURL")

	timeNow := time.Now()

	scenarioMap, err := ParseScenariosFile(scenarioUrl)
	if err != nil { pp.Fatal(err) }

	i:=0
	//for i := 0; i < NUMBEROFCASES; i++ {
	for { // Run forever

		if time.Since(timeNow).Hours() > 5 {
			pp.Println("Updating Scenario file")
			scenarioMap2, err := ParseScenariosFile(scenarioUrl)
			if err != nil { pp.Println("!! Cannot fetch new scenario file, keeping the old one.") 
			} else { scenarioMap = scenarioMap2 }
			timeNow = time.Now()
		}

		// New Visit
		useCase, err := InitUseCase(*debug)
		if err != nil { pp.Fatal(err) }
			
		randScen := 0
		// Random Scenario
		if len(scenarioMap) > 1  { randScen = rand.Intn(len(scenarioMap)) }

		err = ExecuteScenario(scenarioMap[randScen], useCase)
		if err != nil { pp.Fatal(err) }

		// End visit
		useCase.Cua.DeleteVisit()
		time.Sleep(time.Duration(rand.Intn(TIMEBETWEENVISITS)) * time.Second)

		i++
		fmt.Printf("\r%d scenarios executed...", i)
	}
	pp.Println("DONE")
}

// Utils functions
// =====================

func FindDocumentRankByTitle(toFind string, resultList []search.Result) int {
	rank := -1
	for i := 0; i < len(resultList); i++ {
		if strings.Contains(strings.ToLower(resultList[i].Title), strings.ToLower(toFind)) {
			return i
		}
	}
	return rank
}

// Because go "forgot" to include a min function for integers
func Min(a int, b int) int {
	if a < b { return a  
	} else { return b }
}

// Because go "forgot" to include a max function for integers
func Max(a int, b int) int {
	if a > b { return a  
	} else { return b }
}

// Wait a random number of seconds between user actions
func WaitBetweenActions() {
	time.Sleep(time.Duration(rand.Intn(TIMEBETWEENACTIONS)) * time.Second)
}