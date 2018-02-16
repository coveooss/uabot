package scenariolib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/coveo/uabot/defaults"
)

// Config This is the struct that holds all the info on the current bot session.
type Config struct {
	// ScenarioMap A map of the scenarios in the config and their weights.
	ScenarioMap []*Scenario

	// OrgName The name of the Org where you run the bot.
	OrgName string `json:"orgName"`

	// GoodQueries An array of queries that are considered good (return results and good click rank).
	GoodQueries []string `json:"randomGoodQueries"`

	// BadQueries An array of queries that are considered bad (return no results or bad click rank).
	BadQueries []string `json:"randomBadQueries"`

	// GoodQueriesInLang An array of languages containing GoodQueries.
	GoodQueriesInLang map[string][]string `json:"goodQueriesInLanguage"`

	// BadQueriesInLang An array of languages containing BadQueries.
	BadQueriesInLang map[string][]string `json:"badQueriesInLanguage"`

	// Scenarios An array of scenarios to execute
	Scenarios []*Scenario `json:"scenarios"`

	// GlobalFilter A query expression to send along with each queries.
	GlobalFilter string `json:"globalfilter,omitempty"`

	// SearchEndpoint Override of the SearchEndpoint where to send the queries.
	SearchEndpoint string `json:"searchendpoint,omitempty"`

	// AnalyticsEndpoint Override of the default AnalyticsEndpoint where to send analytics.
	AnalyticsEndpoint string `json:"analyticsendpoint,omitempty"`

	// RandomData Override the bot default fake data.
	RandomData *RandomData `json:"randomData,omitempty"`

	// PartialMatch PartialMath param to send with queries.
	PartialMatch bool `json:"partialMatch,omitempty"`

	// PartialMatchKeywords partialMatchKeywords param to send with queries.
	PartialMatchKeywords int `json:"partialMatchKeywords,omitempty"`

	// PartialMatchThreshold partialMatchThreshold param to send with queries.
	PartialMatchThreshold string `json:"partialMatchThreshold,omitempty"`

	// Pipeline The pipeline for the search queries.
	Pipeline string `json:"pipeline,omitempty"`

	// DontWaitBetweenVisits Do not wait between the visits.
	DontWaitBetweenVisits bool `json:"dontWaitBetweenVisits"`

	// DontWaitBetweenActions Do not wait between actions.
	DontWaitBetweenActions bool `json:"dontWaitBetweenActions"`

	// TimeBetweenVisits Time to wait between the visits in seconds
	TimeBetweenVisits int `json:"timeBetweenVisits,omitempty"`

	// TimeBetweenActions The time to wait between actions in seconds
	TimeBetweenActions int `json:"timeBetweenActions,omitempty"`

	// AnonymousThreshold The percentage of visits that are anonymous [0..1].
	AnonymousThreshold float64 `json:"anonymousThreshold,omitempty"`

	// AllowEntitlements Don't use that...
	AllowEntitlements bool `json:"allowEntitlements,omitempty"`

	// RandomCustomData An array of RandomCustomData to send with each queries.
	RandomCustomData []*RandomCustomData `json:"randomCustomData,omitempty"`

	// IsWaitConstant Do you want the wait time to be constant.
	IsWaitConstant bool `json:"isWaitConstant,omitempty"`
}

// RandomData An override of the bot default random/fake data.
type RandomData struct {
	// DefaultOriginLevel1 Override of the default OriginLevel1.
	DefaultOriginLevel1 string `json:"defaultOriginLevel1,omitempty"`

	// DefaultPageViewField Override of the DefaultPageViewField for ALL pageView Events.
	DefaultPageViewField string `json:"defaultPageViewField,omitempty"`

	// Emails Override the defaults fake emails.
	Emails []string `json:"emailSuffixes,omitempty"`

	// FirstNames Override the defaults fake FirstNames.
	FirstNames []string `json:"firstNames,omitempty"`

	// LastNames Override the defaults fake LastNames.
	LastNames []string `json:"lastNames,omitempty"`

	// RandomIPs Override the defaults fake RandomIPs.
	RandomIPs []string `json:"randomIPs,omitempty"`

	// UserAgents Override the defaults fake UserAgents.
	UserAgents []string `json:"useragents,omitempty"`

	// Languages Override the defaults fake Languages.
	Languages []string `json:"languages,omitempty"`

	// MobileUserAgents Override the defaults fake MobileUserAgents.
	MobileUserAgents []string `json:"mobileuseragents, omitempty"`
}

// RandomCustomData Structure of random values for a specific API name.
type RandomCustomData struct {
	APIName string   `json:"apiname"`
	Values  []string `json:"values"`
}

// RandomScenario Returns a random scenario from the list of possible scenarios.
// returns an error if there are no scenarios
func (c *Config) RandomScenario() (*Scenario, error) {
	if len(c.ScenarioMap) < 1 {
		return nil, errors.New("No scenarios detected")
	}
	return c.ScenarioMap[rand.Intn(len(c.ScenarioMap))], nil
}

// RandomQuery Returns a random query good or bad from the list of possible queries.
// returns an error if there are no queries to select from
func (c *Config) RandomQuery(good bool) (string, error) {
	if good {
		if len(c.GoodQueries) < 1 {
			return "", errors.New("No good queries detected")
		}
		return c.GoodQueries[rand.Intn(len(c.GoodQueries))], nil
	}
	if len(c.BadQueries) < 1 {
		return "", errors.New("No bad queries detected")
	}
	return c.BadQueries[rand.Intn(len(c.BadQueries))], nil
}

// RandomQueryInLanguage Returns a random query in a specified language
func (c *Config) RandomQueryInLanguage(good bool, language string) (string, error) {
	if good {
		if len(c.GoodQueriesInLang[language]) < 1 {
			return "", errors.New("No good queries detected in lang : " + language)
		}
		return c.GoodQueriesInLang[language][rand.Intn(len(c.GoodQueriesInLang[language]))], nil
	}
	if len(c.BadQueriesInLang[language]) < 1 {
		return "", errors.New("No bad queries detected in lang : " + language)
	}
	return c.BadQueriesInLang[language][rand.Intn(len(c.BadQueriesInLang[language]))], nil
}

// RandomUserAgent returns a random user agent string to send with an event
func (c *Config) RandomUserAgent(mobile bool) (string, error) {
	if mobile && (len(c.RandomData.MobileUserAgents) > 0) {
		return c.RandomData.MobileUserAgents[rand.Intn(len(c.RandomData.MobileUserAgents))], nil
	} else if len(c.RandomData.UserAgents) > 0 || len(c.RandomData.MobileUserAgents) > 0 {
		agents := append(c.RandomData.UserAgents, c.RandomData.MobileUserAgents...)
		return agents[rand.Intn(len(agents))], nil
	}
	return "", errors.New("Cannot find any user agents")
}

// RandomIP returns a random IP to send with an event
func (c *Config) RandomIP() (string, error) {
	if len(c.RandomData.RandomIPs) < 1 {
		return "", errors.New("Cannot find any random IP")
	}

	return c.RandomData.RandomIPs[rand.Intn(len(c.RandomData.RandomIPs))], nil
}

// NewConfigFromPath Create a new config from a JSON config file path
//
// jsonPath The path to the JSON file.
func NewConfigFromPath(jsonPath string) (*Config, error) {

	file, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading JSON file : %v", err)
	}

	c := &Config{}
	err = json.Unmarshal(file, c)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON : %v", err)
	}

	fillDefaults(c)

	err = c.makeScenarioMap()
	if err != nil {
		return nil, fmt.Errorf("Error making scenario map : %v", err)
	}
	return c, nil
}

// NewConfigFromURL Create a new Config from an url to a JSON file
//
// jsonURL The URL to the json config path
func NewConfigFromURL(jsonURL string) (*Config, error) {

	resp, err := http.Get(jsonURL)
	if err != nil {
		return nil, errors.New("Cannot read JSON config file")
	}
	defer resp.Body.Close()

	c := &Config{}
	err = json.NewDecoder(resp.Body).Decode(c)
	if err != nil {
		return nil, errors.New("Cannot decode JSON file")
	}

	fillDefaults(c)

	err = c.makeScenarioMap()
	if err != nil {
		return nil, errors.New("Cannot make the scenario map")
	}
	return c, nil
}

// makeScenarioMap Private function to create the map of scenarios
// from the config that was built from a json file
func (c *Config) makeScenarioMap() error {
	scenarioMap := []*Scenario{}
	totalWeight := 0
	for i := 0; i < len(c.Scenarios); i++ {
		weight := c.Scenarios[i].Weight
		totalWeight += weight
		for j := 0; j < weight; j++ {
			scenarioMap = append(scenarioMap, c.Scenarios[i])
		}
	}
	c.ScenarioMap = scenarioMap
	return nil
}

func fillDefaults(c *Config) {
	if len(c.RandomData.FirstNames) == 0 {
		c.RandomData.FirstNames = defaults.FIRSTNAMES
	}

	if len(c.RandomData.LastNames) == 0 {
		c.RandomData.LastNames = defaults.LASTNAMES
	}

	if len(c.RandomData.Emails) == 0 {
		c.RandomData.Emails = defaults.EMAILS
	}

	if len(c.RandomData.RandomIPs) == 0 {
		c.RandomData.RandomIPs = defaults.IPS
	}

	if len(c.RandomData.UserAgents) == 0 {
		c.RandomData.UserAgents = defaults.USERAGENTS
	}

	if len(c.RandomData.MobileUserAgents) == 0 {
		c.RandomData.MobileUserAgents = defaults.MOBILEUSERAGENTS
	}

	if c.SearchEndpoint == "" {
		c.SearchEndpoint = defaults.SEARCHENDPOINT_PROD
	}

	if c.AnalyticsEndpoint == "" {
		c.AnalyticsEndpoint = defaults.ANALYTICSENDPOINT_PROD
	}

	if c.RandomData.DefaultPageViewField == "" {
		c.RandomData.DefaultPageViewField = defaults.DEFAULTPAGEVIEWFIELD
	}

	if c.RandomData.DefaultOriginLevel1 == "" {
		c.RandomData.DefaultOriginLevel1 = defaults.DEFAULTORIGIN1
	}
}
