package visit

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"

	"github.com/erocheleau/uabot/defaults"
)

// BotConfig All the information necessary to run the ua bot
type BotConfig struct {
	OrgName               string              `json:"orgName"`
	Queries               QueriesDataSet      `json:"queriesDataSet"`
	Scenarios             []*Scenario         `json:"scenarios"`
	SearchEndpoint        string              `json:"searchEndpoint,omitempty"`
	AnalyticsEndpoint     string              `json:"analyticsEndpoint,omitempty"`
	Users                 UserDataSet         `json:"userDataSet,omitempty"`
	TimeBetweenVisits     int                 `json:"timeBetweenVisits,omitempty"`
	TimeBetweenActions    int                 `json:"timeBetweenActions,omitempty"`
	AllowAnonymous        bool                `json:"allowAnonymousVisits,omitempty"`
	AnonymousTreshold     float64             `json:"anonymousTreshold,omitempty"`
	AllowEntitlements     bool                `json:"allowEntitlements,omitempty"`
	RandomCustomData      []*RandomCustomData `json:"randomCustomData,omitempty"`
	RandomDocumentAuthors []string            `json:"randomAuthors,omitempty"`
	ScenarioMap           ScenarioMap
}

type QueryParams struct {
	PartialMatch          bool   `json:"partialMatch,omitempty"`
	PartialMatchKeywords  int    `json:"partialMatchKeywords,omitempty"`
	PartialMatchThreshold string `json:"partialMatchThreshold,omitempty"`
	Pipeline              string `json:"pipeline,omitempty"`
	DefaultOriginLevel1   string `json:"defaultOriginLevel1,omitempty"`
	GlobalFilter          string `json:"globalfilter,omitempty"`
}

// QueriesDataSet The dataset of random queries that the bot can use
type QueriesDataSet struct {
	GoodQueries []string `json:"goodQueries"`
	BadQueries  []string `json:"badQueries"`
}

// UserDataSet The dataset of random user information the bot can use
type UserDataSet struct {
	Emails           []string `json:"emailSuffixes,omitempty"`
	FirstNames       []string `json:"firstNames,omitempty"`
	LastNames        []string `json:"lastNames,omitempty"`
	RandomIPs        []string `json:"randomIPs,omitempty"`
	UserAgents       []string `json:"useragents,omitempty"`
	MobileUserAgents []string `json:"mobileuseragents, omitempty"`
	Languages        []string `json:"languages,omitempty"`
}

type RandomCustomData struct {
	APIName string   `json:"apiname"`
	Values  []string `json:"values"`
}

// NewConfigFromURL Create a new BotConfig from an hosted JSON file
func NewConfigFromURL(jsonURL string) (*BotConfig, error) {
	resp, err := http.Get(jsonURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	c := &BotConfig{}
	err = json.NewDecoder(resp.Body).Decode(c)
	if err != nil {
		return nil, err
	}

	fillDefaults(c)

	err = c.makeScenarioMap()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// makeScenarioMap Private function to create the map of scenarios
// from the config that was built from a json file
func (c *BotConfig) makeScenarioMap() error {
	if len(c.Scenarios) < 1 {
		return errors.New("No scenarios found!")
	}
	var scenarioMap = ScenarioMap{}
	totalWeight := 0
	iter := 0
	for i := 0; i < len(c.Scenarios); i++ {
		weight := c.Scenarios[i].Weight
		totalWeight += weight
		for j := 0; j < weight; j++ {
			scenarioMap[iter] = c.Scenarios[i]
			iter++
		}
	}
	c.ScenarioMap = scenarioMap
	return nil
}

// RandomScenario Returns a random scenario from the list of possible scenarios.
// returns an error if there are no scenarios
func (c *BotConfig) RandomScenario() (*Scenario, error) {
	if len(c.ScenarioMap) < 1 {
		return nil, errors.New("No scenarios detected")
	}
	return c.ScenarioMap[rand.Intn(len(c.ScenarioMap))], nil
}

func fillDefaults(c *BotConfig) {
	if len(c.Users.FirstNames) == 0 {
		c.Users.FirstNames = defaults.FIRSTNAMES
	}
	if len(c.Users.LastNames) == 0 {
		c.Users.LastNames = defaults.LASTNAMES
	}
	if len(c.Users.Emails) == 0 {
		c.Users.Emails = defaults.EMAILS
	}
	if len(c.Users.RandomIPs) == 0 {
		c.Users.RandomIPs = defaults.IPS
	}
	if len(c.Users.UserAgents) == 0 {
		c.Users.UserAgents = defaults.USERAGENTS
	}
	if len(c.Users.MobileUserAgents) == 0 {
		c.Users.MobileUserAgents = defaults.MOBILEUSERAGENTS
	}
	if c.SearchEndpoint == "" {
		c.SearchEndpoint = defaults.SEARCHENDPOINT_PROD
	}
	if c.AnalyticsEndpoint == "" {
		c.AnalyticsEndpoint = defaults.ANALYTICSENDPOINT_PROD
	}
	if len(c.RandomDocumentAuthors) == 0 {
		c.RandomDocumentAuthors = defaults.AUTHORNAMES
	}
}
