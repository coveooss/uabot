package main

import (
	"log"
	"os"
	"flag"
	"strings"
	"github.com/adambbolduc/uabot/explorerlib"
	"github.com/erocheleau/uabot/scenariolib"
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
	PARAM_FETCH_NUMBER_OF_RESULTS = flag.Int("fetchQueryNumber", 100, "Amount of results for each query")
	PARAM_DOCUMENTS_EXPLORATION_PERCENTAGE = flag.Float64("explorationRatio", 0.5, "Query Response Result Considered / Amount of Documents in Index")
	PARAM_FIELD_TO_EXPLORE_EQUALLY = flag.String("fields", "", "Fields to Explore Equally, seperated by comas ex: @syssource,@filetype")
	PARAM_INDEX_ENDPOINT = flag.String("searchEndpoint", "", "Search Endpoint Url")
	PARAM_SEARCH_TOKEN = flag.String("token", "", "Search Token for the organization")
	PARAM_MAXIMUM_NUMBER_OF_QUERY_BY_LANGUAGE = flag.Int("maxNumberOfQueryPerLanguage", 10, "The Maximum number of query expression for each languages")
	PARAM_ANALYTICS_ENDPOINT = flag.String("analyticsEndpoint", "", "Analytics Endpoint Url")
	PARAM_ORG_NAME = flag.String("org", "", "Name of the Organization")
	PARAM_OUTPUT_FILEPATH = flag.String("outputFilepath" , "out.json", "Output Filename to dump JSON Config File for uabot")
)

func main() {
	flag.Parse()
	Info = log.New(os.Stdout, "INFO | ", log.Ldate | log.Ltime)
	index, status := explorerlib.NewIndex(*PARAM_INDEX_ENDPOINT, *PARAM_SEARCH_TOKEN)
	check(status)

	wordCountsByLanguage, status := explorerlib.FindWordsByLanguageInIndex(
		index,
		strings.Split(*PARAM_FIELD_TO_EXPLORE_EQUALLY, ","),
		*PARAM_DOCUMENTS_EXPLORATION_PERCENTAGE,
		*PARAM_FETCH_NUMBER_OF_RESULTS)
	check(status)

	languages, status := index.FetchLanguages();
	check(status)

	err := explorerlib.NewScenarioBuilder().WithOrgName(*PARAM_ORG_NAME).WithSearchEndpoint(*PARAM_INDEX_ENDPOINT).WithAnalyticsEndpoint(*PARAM_ANALYTICS_ENDPOINT).AllAnonymous().WithLanguages(languages).WithWordCountsByLanguage(wordCountsByLanguage, *PARAM_MAXIMUM_NUMBER_OF_QUERY_BY_LANGUAGE, 2).WithTimeBetweenActions(1).WithTimeBetweenVisits(5).WithScenarios(&scenariolib.Scenario{
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
					"probability": 1.0,
				},
			},
		},
		Weight: 1,
	}).Save(*PARAM_OUTPUT_FILEPATH)
	check(err)
}
