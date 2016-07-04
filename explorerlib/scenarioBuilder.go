package explorerlib

import (
	"github.com/erocheleau/uabot/scenariolib"
	"encoding/json"
	"io/ioutil"
)

type ScenarioBuilder interface {
	build(*scenariolib.Config)
}

type scenarioBuilder struct {
	config scenariolib.Config
}

func NewScenarioBuilder() (*scenarioBuilder) {
	return &scenarioBuilder{
		config:scenariolib.Config{},
	}
}
func (builder *scenarioBuilder) WithOrgName(orgName string) *scenarioBuilder {
	builder.config.OrgName = orgName
	return builder
}

func (builder *scenarioBuilder) WithLanguages(languages []string) *scenarioBuilder {
	builder.config.Languages = languages
	return builder
}

func (builder *scenarioBuilder) AllAnonymous() *scenarioBuilder {
	builder.config.AllowAnonymous = true
	builder.config.AnonymousThreshold = 1
	return builder
}

func (builder *scenarioBuilder) WithSearchEndpoint(endpoint string) *scenarioBuilder {
	builder.config.SearchEndpoint = endpoint
	return builder
}

func (builder *scenarioBuilder) WithAnalyticsEndpoint(endpoint string) *scenarioBuilder {
	builder.config.AnalyticsEndpoint = endpoint
	return builder
}

func (builder *scenarioBuilder) WithTimeBetweenActions(time int) *scenarioBuilder {
	builder.config.TimeBetweenActions = time
	return builder
}

func (builder *scenarioBuilder) WithTimeBetweenVisits(time int) *scenarioBuilder {
	builder.config.TimeBetweenVisits = time
	return builder
}

func (builder *scenarioBuilder) WithGoodQueryByLanguage(goodQueriesByLanguage map[string][]string) *scenarioBuilder {
	builder.config.GoodQueriesInLang = goodQueriesByLanguage
	return builder
}

func (builder *scenarioBuilder) WithScenarios(scenarios []*scenariolib.Scenario) *scenarioBuilder {
	builder.config.Scenarios = scenarios
	return builder
}

func (builder *scenarioBuilder) Build() (*scenariolib.Config) {
	return &scenariolib.Config{}
}

func (builder *scenarioBuilder) Save(path string) error {
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


