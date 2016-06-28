package main

import (
	"log"
	"os"
	"github.com/adambbolduc/uabot/natural-language"
	"github.com/erocheleau/uabot/scenariolib"
	"io/ioutil"
	"encoding/json"
	"flag"
	"strings"
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
)

func main() {
	flag.Parse()
	Info = log.New(os.Stdout, "INFO | ", log.Ldate | log.Ltime)
	index, status := natural_language.NewIndex(*PARAM_INDEX_ENDPOINT, *PARAM_SEARCH_TOKEN)
	check(status)

	fields := strings.Split(*PARAM_FIELD_TO_EXPLORE_EQUALLY, ",")

	wordCountsByLanguage, status := natural_language.FindWordsByLanguageInIndex(index, fields, *PARAM_DOCUMENTS_EXPLORATION_PERCENTAGE, *PARAM_FETCH_NUMBER_OF_RESULTS)
	check(status)

	queriesInLanguage := make(map[string][]string)
	for language, wordCounts := range wordCountsByLanguage {
		words := []string{}
		for i, word := range wordCounts.Words {
			if i < *PARAM_MAXIMUM_NUMBER_OF_QUERY_BY_LANGUAGE{
				words = append(words, word.Word)
			} else {
				break
			}
		}
		queriesInLanguage[language] = words
	}
	config := scenariolib.Config{}
	config.GoodQueriesInLang = queriesInLanguage
	bytes, err := json.Marshal(config)
	check(err)
	writeerr := ioutil.WriteFile("out.json", bytes, 0644)
	check(writeerr)
}
