package main

import (
	"log"
	"os"
	"github.com/adambbolduc/uabot/explorerlib"
	"github.com/erocheleau/uabot/scenariolib"
	"encoding/json"
	"io/ioutil"
	"flag"
)

var (
	Info *log.Logger
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var (
	configFile = flag.String("config", "config.json", "File Path for configuration File")
)

func main() {
	Info = log.New(os.Stdout, "INFO | ", log.Ldate | log.Ltime)

	var config explorerlib.Config
	bytes, readErr := ioutil.ReadFile(*configFile)
	check(readErr)
	json.Unmarshal(bytes, &config)

	index, status := explorerlib.NewIndex(config.SearchEndpoint, config.SearchToken)
	check(status)

	wordCountsByLanguage, status := explorerlib.FindWordsByLanguageInIndex(
		index,
		config.FieldsToExploreEqually,
		config.DocumentsExplorationPercentage,
		config.FetchNumberOfResults)
	check(status)

	languages, status := index.FetchLanguages();
	check(status)

	err := explorerlib.NewScenarioBuilder().WithOrgName(config.Org).WithSearchEndpoint(config.SearchEndpoint).WithAnalyticsEndpoint(config.AnalyticsEndpoint).AllAnonymous().WithLanguages(languages).WithWordCountsByLanguage(wordCountsByLanguage, config.NumberOfQueryByLanguage, config.AverageNumberOfWordsPerQuery).WithTimeBetweenActions(1).WithTimeBetweenVisits(5).WithScenarios(&scenariolib.Scenario{
		Name: "scenario1",
		Events:[]scenariolib.JSONEvent{
			scenariolib.JSONEvent{
				Type:"Search",
				Arguments: map[string]interface{}{
					"queryText" : "",
					"goodQuery" : true,
					"matchLanguage":true,
				},
			},
			scenariolib.JSONEvent{
				Type:"Click",
				Arguments: map[string]interface{}{
					"docNo": -1,
					"offset": 0,
					"probability": 0.3,
				},
			},
			scenariolib.JSONEvent{
				Type:"View",
				Arguments: map[string]interface{}{
					"docNo": -1,
					"offset": 0,
					"probability": 0.3,
				},
			},
			scenariolib.JSONEvent{
				Type: "Search",
				Arguments: map[string]interface{}{
					"queryText" :"",
					"goodQuery":true,
					"matchLanguage":true,
				},
			},
			scenariolib.JSONEvent{
				Type:"Click",
				Arguments: map[string]interface{}{
					"docNo": 1,
					"offset": 0,
					"probability": 0.5,
				},
			},
			scenariolib.JSONEvent{
				Type:"View",
				Arguments: map[string]interface{}{
					"docNo": 1,
					"offset": 0,
					"probability": 0.5,
				},
			},
		},
		Weight: 1,
	}).Save(config.OutputFilePath)
	check(err)
}
