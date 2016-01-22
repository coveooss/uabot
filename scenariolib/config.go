package scenariolib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

var _scenariosMap = map[int]*Scenario{}

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
	SearchEndpoint    string      `json:"searchendpoint"`
	AnalyticsEndpoint string      `json:"analyticsendpoint"`
	OrgName           string      `json:"orgName"`
	Emails            []string    `json:"emailSuffixes"`
	FirstNames        []string    `json:"firstNames"`
	LastNames         []string    `json:"lastNames"`
	GoodQueries       []string    `json:"randomGoodQueries"`
	BadQueries        []string    `json:"randomBadQueries"`
	RandomIPs         []string    `json:"randomIPs"`
	Scenarios         []*Scenario `json:"scenarios"`
	ScenarioMap       map[int]*Scenario
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
	c := Config{}
	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		return nil, errors.New("Cannot decode JSON file")
	}
	err = c.makeScenarioMap()
	if err != nil {
		return nil, errors.New("Cannot make the scenario map")
	}
	return &c, nil
}

// makeScenarioMap Private function to create the map of scenarios
// from the config that was built from a json file
func (c *Config) makeScenarioMap() error {
	scenarioMap := map[int]*Scenario{}
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
