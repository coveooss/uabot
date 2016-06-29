package main

import (
	"os"
	"github.com/adambbolduc/uabot/explorerlib"
	"github.com/erocheleau/uabot/scenariolib"
	"encoding/json"
	"io/ioutil"
	"flag"
	"time"
	"math/rand"
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

	languages, status := index.FetchLanguages();
	check(status)

	err := explorerlib.NewScenarioBuilder().WithOrgName(config.Org).WithSearchEndpoint(config.SearchEndpoint).WithAnalyticsEndpoint(config.AnalyticsEndpoint).AllAnonymous().WithLanguages(languages).WithWordCountsByLanguage(wordCountsByLanguage, config.NumberOfQueryByLanguage, config.AverageNumberOfWordsPerQuery).WithTimeBetweenActions(1).WithTimeBetweenVisits(5).WithScenarios(config.MainScenario).Save(config.OutputFilePath)
	check(err)

	scenariolib.Info.Println("Running Bot")

	conf, err := scenariolib.NewConfigFromPath(config.OutputFilePath)
	check(err)

	uabot := scenariolib.NewUabot(true, conf, config.OutputFilePath, config.SearchToken, config.AnalyticsToken, random)
	err =uabot.Run()
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

