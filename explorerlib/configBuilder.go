package explorerlib

import (
	"github.com/erocheleau/uabot/scenariolib"
	"encoding/json"
	"io/ioutil"
)

type botConfigurationBuilder struct {
	config scenariolib.Config
}

func NewBotConfigurationBuilder() (*botConfigurationBuilder) {
	return &botConfigurationBuilder{
		config:scenariolib.Config{},
	}
}

func (builder *botConfigurationBuilder) WithOrgName(orgName string) *botConfigurationBuilder {
	builder.config.OrgName = orgName
	return builder
}

func (builder *botConfigurationBuilder) WithLanguages(languages []string) *botConfigurationBuilder {
	builder.config.Languages = languages
	return builder
}

func (builder *botConfigurationBuilder) AllAnonymous() *botConfigurationBuilder {
	builder.config.AllowAnonymous = true
	builder.config.AnonymousThreshold = 1
	return builder
}

func (builder *botConfigurationBuilder) WithSearchEndpoint(endpoint string) *botConfigurationBuilder {
	builder.config.SearchEndpoint = endpoint
	return builder
}

func (builder *botConfigurationBuilder) WithAnalyticsEndpoint(endpoint string) *botConfigurationBuilder {
	builder.config.AnalyticsEndpoint = endpoint
	return builder
}

func (builder *botConfigurationBuilder) WithTimeBetweenActions(time int) *botConfigurationBuilder {
	builder.config.TimeBetweenActions = time
	return builder
}

func (builder *botConfigurationBuilder) WithTimeBetweenVisits(time int) *botConfigurationBuilder {
	builder.config.TimeBetweenVisits = time
	return builder
}

func (builder *botConfigurationBuilder) WithGoodQueryByLanguage(goodQueriesByLanguage map[string][]string) *botConfigurationBuilder {
	builder.config.GoodQueriesInLang = goodQueriesByLanguage
	return builder
}

func (builder *botConfigurationBuilder) WithScenarios(scenarios []*scenariolib.Scenario) *botConfigurationBuilder{
	builder.config.Scenarios = scenarios
	return builder
}

func (builder *botConfigurationBuilder) NoWait() *botConfigurationBuilder {
	builder.config.DontWaitBetweenVisits = true
	builder.config.DontWaitBetweenActions = true
	return builder
}

func (builder *botConfigurationBuilder) Build() (*scenariolib.Config) {
	return &scenariolib.Config{}
}

func (builder *botConfigurationBuilder) Save(path string) error {
	bytes, err := json.Marshal(builder.config)
	if err != nil {
		return err
	}
	writeerr := ioutil.WriteFile(path, bytes, 0644)
	if writeerr != nil {
		return writeerr
	}
	return nil
}


