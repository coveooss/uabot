package autobot

import (
	"github.com/adambbolduc/uabot/explorerlib"
	"github.com/erocheleau/uabot/scenariolib"
	"math/rand"
)

type Autobot struct {
	config *explorerlib.Config
	random *rand.Rand
}

func NewAutobot(_config *explorerlib.Config, _random *rand.Rand) *Autobot {
	return &Autobot{
		config: _config,
		random: _random,
	}
}

func (bot *Autobot) Run(quitChannel chan bool) error {
	index, status := explorerlib.NewIndex(bot.config.SearchEndpoint, bot.config.SearchToken)

	wordCountsByLanguage, status := explorerlib.FindWordsByLanguageInIndex(
		index,
		bot.config.FieldsToExploreEqually,
		bot.config.DocumentsExplorationPercentage,
		bot.config.FetchNumberOfResults)
	if status != nil {
		return status
	}

	languages, status := index.Client.ListFacetValues("@language", 1000)
	if status != nil {
		return status
	}

	goodQueries, status := index.BuildGoodQueries(wordCountsByLanguage, bot.config.NumberOfQueryByLanguage, bot.config.AverageNumberOfWordsPerQuery)
	if status != nil {
		return status
	}

	taggedLanguages := make([]string, 0)
	scenarios := []*scenariolib.Scenario{}

	for _, lang := range languages.Values {
		taggedLanguage := explorerlib.LanguageToTag(lang.Value)
		taggedLanguages = append(taggedLanguages, taggedLanguage)
		scenario := explorerlib.NewScenarioBuilder().WithName("search and click in " + lang.Value).WithWeight(lang.NumberOfResults).WithLanguage(taggedLanguage).WithEvent(explorerlib.NewSearchEvent(true)).WithEvent(explorerlib.NewClickEvent(0.4)).WithEvent(explorerlib.NewSearchEvent(true)).WithEvent(explorerlib.NewClickEvent(0.8)).Build()
		scenarios = append(scenarios, scenario)
		viewScenarioBuilder := explorerlib.NewScenarioBuilder().WithName("views in " + lang.Value).WithWeight(lang.NumberOfResults).WithLanguage(taggedLanguage).WithEvent(explorerlib.NewSearchEvent(false))
		for i := 0; i < 20; i++ {
			viewScenarioBuilder.WithEvent(explorerlib.NewViewEvent())
		}
		scenarios = append(scenarios, viewScenarioBuilder.Build())
	}

	err := explorerlib.NewBotConfigurationBuilder().WithOrgName(bot.config.Org).WithSearchEndpoint(bot.config.SearchEndpoint).WithAnalyticsEndpoint(bot.config.AnalyticsEndpoint).AllAnonymous().WithLanguages(taggedLanguages).WithGoodQueryByLanguage(goodQueries).WithTimeBetweenActions(1).WithTimeBetweenVisits(5).WithScenarios(scenarios).NoWait().Save(bot.config.OutputFilePath)
	if err != nil {
		return err
	}

	scenariolib.Info.Println("Running Bot")

	uabot := scenariolib.NewUabot(true, bot.config.OutputFilePath, bot.config.SearchToken, bot.config.AnalyticsToken, bot.random)
	err = uabot.Run(quitChannel)
	return err
}

func (bot *Autobot) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"searchEndpoint":                 bot.config.SearchEndpoint,
		"analyticsEndpoint":              bot.config.AnalyticsEndpoint,
		"averageNumberOfWordsPerQuery":   bot.config.AverageNumberOfWordsPerQuery,
		"documentsExplorationPercentage": bot.config.DocumentsExplorationPercentage,
		"fieldsToExploreEqually":         bot.config.FieldsToExploreEqually,
		"org":                      bot.config.Org,
		"outputFilepath":           bot.config.OutputFilePath,
		"numberOfQueryPerLanguage": bot.config.NumberOfQueryByLanguage,
		"numberOfResultsPerQuery":  bot.config.FetchNumberOfResults,
	}
}
