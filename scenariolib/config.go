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
//
// OrgName     The name of the Org where you run the bot.
// Emails      An array of strings containing email suffixes to use to randomize name of visitors.
// FirstNames  An array of strings containing first names to use to randomize name of visitors.
// LastNames   An array of strings containing last names to use to randomize name of visitors.
// GoodQueries An array of queries that are considered good (return results and good click rank).
// BadQueries  An array of queries that are considered bad (return no results or bad click rank).
// RandomIPs   An array of IPs to randomize the location of visits
// Scenarios   An array of scenarios that can happen.
// ScenarioMap A map that will be built with the scenarios and their respective weights.
type Config struct {
	ScenarioMap            []*Scenario
	OrgName                string              `json:"orgName"`
	GoodQueries            []string            `json:"randomGoodQueries"`
	BadQueries             []string            `json:"randomBadQueries"`
	GoodQueriesInLang      map[string][]string `json:"goodQueriesInLanguage"`
	BadQueriesInLang       map[string][]string `json:"badQueriesInLanguage"`
	Scenarios              []*Scenario         `json:"scenarios"`
	DefaultOriginLevel1    string              `json:"defaultOriginLevel1,omitempty"`
	DefaultPageViewField   string              `json:"defaultPageViewField,omitempty"`
	GlobalFilter           string              `json:"globalfilter,omitempty"`
	SearchEndpoint         string              `json:"searchendpoint,omitempty"`
	AnalyticsEndpoint      string              `json:"analyticsendpoint,omitempty"`
	Emails                 []string            `json:"emailSuffixes,omitempty"`
	FirstNames             []string            `json:"firstNames,omitempty"`
	LastNames              []string            `json:"lastNames,omitempty"`
	RandomIPs              []string            `json:"randomIPs,omitempty"`
	UserAgents             []string            `json:"useragents,omitempty"`
	Languages              []string            `json:"languages,omitempty"`
	MobileUserAgents       []string            `json:"mobileuseragents, omitempty"`
	PartialMatch           bool                `json:"partialMatch,omitempty"`
	PartialMatchKeywords   int                 `json:"partialMatchKeywords,omitempty"`
	PartialMatchThreshold  string              `json:"partialMatchThreshold,omitempty"`
	Pipeline               string              `json:"pipeline,omitempty"`
	DontWaitBetweenVisits  bool                `json:"dontWaitBetweenVisits"`
	DontWaitBetweenActions bool                `json:"dontWaitBetweenActions"`
	TimeBetweenVisits      int                 `json:"timeBetweenVisits,omitempty"`
	TimeBetweenActions     int                 `json:"timeBetweenActions,omitempty"`
	AllowAnonymous         bool                `json:"allowAnonymousVisits,omitempty"`
	AnonymousThreshold     float64             `json:"anonymousThreshold,omitempty"`
	AllowEntitlements      bool                `json:"allowEntitlements,omitempty"`
	RandomCustomData       []*RandomCustomData `json:"randomCustomData,omitempty"`
}

type RandomCustomData struct {
	APIName string   `json:"apiname"`
	Values  []string `json:"values"`
}

// DEFAULTANONYMOUSTHRESHOLD The default portion of users who are anonymous
const DEFAULTANONYMOUSTHRESHOLD float64 = 0.5

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

// Returns a random query in a specified language
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
	if mobile && (len(c.MobileUserAgents) > 0) {
		return c.MobileUserAgents[rand.Intn(len(c.MobileUserAgents))], nil
	} else if len(c.UserAgents) > 0 || len(c.MobileUserAgents) > 0 {
		agents := append(c.UserAgents, c.MobileUserAgents...)
		return agents[rand.Intn(len(agents))], nil
	}
	return "", errors.New("Cannot find any user agents")
}

// RandomIP returns a random IP to send with an event
func (c *Config) RandomIP() (string, error) {
	if len(c.RandomIPs) < 1 {
		return "", errors.New("Cannot find any random IP")
	}

	return c.RandomIPs[rand.Intn(len(c.RandomIPs))], nil
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
	if len(c.FirstNames) == 0 {
		c.FirstNames = defaults.FIRSTNAMES
	}

	if len(c.LastNames) == 0 {
		c.LastNames = defaults.LASTNAMES
	}

	if len(c.Emails) == 0 {
		c.Emails = defaults.EMAILS
	}

	if len(c.RandomIPs) == 0 {
		c.RandomIPs = defaults.IPS
	}

	if len(c.UserAgents) == 0 {
		c.UserAgents = defaults.USERAGENTS
	}

	if len(c.MobileUserAgents) == 0 {
		c.MobileUserAgents = defaults.MOBILEUSERAGENTS
	}

	if c.SearchEndpoint == "" {
		c.SearchEndpoint = defaults.SEARCHENDPOINT_PROD
	}

	if c.AnalyticsEndpoint == "" {
		c.AnalyticsEndpoint = defaults.ANALYTICSENDPOINT_PROD
	}
}
