package autobot

import (
	"github.com/erocheleau/uabot/explorerlib"
	"github.com/erocheleau/uabot/scenariolib"
	"math/rand"
	"strconv"
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

	for _, lang := range languages.Values { // For each language
		taggedLanguage := explorerlib.LanguageToTag(lang.Value)
		taggedLanguages = append(taggedLanguages, taggedLanguage)
		for originLevel1, originLevels2 := range bot.config.OriginLevels {
			for _, originLevel2 := range originLevels2 {
				for i := 1; i <= 5; i++ {
					scenarioBuilder := explorerlib.NewScenarioBuilder().WithName(strconv.Itoa(i) + " search and click in " + lang.Value).WithWeight(lang.NumberOfResults).WithLanguage(taggedLanguage).WithEvent(explorerlib.NewSetOriginEvent(originLevel1, originLevel2))
					for j := 1; j <= i; j++ {
						scenarioBuilder = scenarioBuilder.WithEvent(explorerlib.NewSearchEvent(true)).WithEvent(explorerlib.NewClickEvent(0.4)).WithEvent(explorerlib.NewClickEvent(0.8))
					}
					scenarios = append(scenarios, scenarioBuilder.Build())
				}

				viewScenarioBuilder := explorerlib.NewScenarioBuilder().WithName("views in " + lang.Value).WithWeight(lang.NumberOfResults).WithLanguage(taggedLanguage).WithEvent(explorerlib.NewSetOriginEvent(originLevel1, originLevel2)).WithEvent(explorerlib.NewSearchEvent(false))
				for i := 0; i < 20; i++ {
					viewScenarioBuilder.WithEvent(explorerlib.NewViewEvent())
				}
				scenarios = append(scenarios, viewScenarioBuilder.Build())
			}
		}

	}

	err := explorerlib.NewBotConfigurationBuilder().WithOrgName(bot.config.Org).WithSearchEndpoint(bot.config.SearchEndpoint).WithAnalyticsEndpoint(bot.config.AnalyticsEndpoint).AllAnonymous().WithLanguages(taggedLanguages).WithGoodQueryByLanguage(goodQueries).WithTimeBetweenActions(1).WithTimeBetweenVisits(5).WithScenarios(scenarios).NoWait().WithOriginLevels(bot.config.OriginLevels).Save(bot.config.OutputFilePath)
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
		"originLevels":             bot.config.OriginLevels,
	}
}
