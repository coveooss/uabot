package main

import (
	"encoding/json"
	"flag"
	"github.com/adambbolduc/uabot/explorerlib"
	"github.com/erocheleau/uabot/scenariolib"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

var (
	configFile = flag.String("config", "config.json", "File Path for configuration File")
)

func main() {
	scenariolib.InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	source := rand.NewSource(int64(time.Now().Unix()))
	random := rand.New(source)

	scenariolib.Info.Println("Reading config file")
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

	languages, status := index.Client.ListFacetValues("@language",1000)
	check(status)

	goodQueries, status := index.BuildGoodQueries(wordCountsByLanguage, config.NumberOfQueryByLanguage, config.AverageNumberOfWordsPerQuery)
	check(status)

	taggedLanguages := make([]string, 0)
	scenarios := []*scenariolib.Scenario{}

	for _, lang := range languages.Values {
		taggedLanguage := explorerlib.LanguageToTag(lang.Value)
		taggedLanguages = append(taggedLanguages, taggedLanguage)
		scenario := explorerlib.NewScenarioBuilder().WithName("search and click in " + lang.Value).WithWeight(lang.NumberOfResults).WithLanguage(taggedLanguage).WithEvent(explorerlib.NewSearchEvent(true)).WithEvent(explorerlib.NewClickEvent(0.4)).WithEvent(explorerlib.NewSearchEvent(true)).WithEvent(explorerlib.NewClickEvent(0.8)).Build()
		scenarios = append(scenarios, scenario)
		viewScenarioBuilder := explorerlib.NewScenarioBuilder().WithName("views in " + lang.Value).WithWeight(lang.NumberOfResults).WithLanguage(taggedLanguage).WithEvent(explorerlib.NewSearchEvent(false))
		for i := 0 ; i < 20 ; i++ {
			viewScenarioBuilder.WithEvent(explorerlib.NewViewEvent())
		}
		scenarios = append(scenarios, viewScenarioBuilder.Build())
	}

	err := explorerlib.NewBotConfigurationBuilder().WithOrgName(config.Org).WithSearchEndpoint(config.SearchEndpoint).WithAnalyticsEndpoint(config.AnalyticsEndpoint).AllAnonymous().WithLanguages(taggedLanguages).WithGoodQueryByLanguage(goodQueries).WithTimeBetweenActions(1).WithTimeBetweenVisits(5).WithScenarios(scenarios).Save(config.OutputFilePath)
	check(err)

	scenariolib.Info.Println("Running Bot")

	uabot := scenariolib.NewUabot(true, config.OutputFilePath, config.SearchToken, config.AnalyticsToken, random)
	err = uabot.Run()
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
