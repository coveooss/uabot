package explorerlib

import (
	"github.com/erocheleau/uabot/scenariolib"
)

type scenarioBuilder struct {
	scenario scenariolib.Scenario
}

func (builder *scenarioBuilder) WithLanguage(language string) *scenarioBuilder {
	builder.scenario.Language = language
	return builder
}

func (builder *scenarioBuilder) WithName(name string) *scenarioBuilder {
	builder.scenario.Name = name
	return builder
}

func (builder *scenarioBuilder) WithUserAgent(userAgent string) *scenarioBuilder {
	builder.scenario.UserAgent = userAgent
	return builder
}

func (builder *scenarioBuilder) WithWeight(weight int) *scenarioBuilder {
	builder.scenario.Weight = weight
	return builder
}

func (builder *scenarioBuilder) WithEvent(event scenariolib.JSONEvent) *scenarioBuilder {
	builder.scenario.Events = append(builder.scenario.Events, event)
	return builder
}

func (builder *scenarioBuilder) Build() *scenariolib.Scenario {
	return &builder.scenario
}

func NewScenarioBuilder() *scenarioBuilder {
	return &scenarioBuilder{}
}

func NewSearchEvent(log bool) scenariolib.JSONEvent {
	return scenariolib.JSONEvent{
		Type: "Search",
		Arguments: map[string]interface{}{
			"queryText": "",
			"logEvent": log,
			"goodQuery": true,
			"matchLanguage": true,
			"caseSearch": false,
		},
	}
}

func NewClickEvent(probability float64) scenariolib.JSONEvent {
	return scenariolib.JSONEvent{
		Type: "Click",
		Arguments: map[string]interface{}{
			"offset": 0,
			"probability": probability,
			"docNo": -1,
		},
	}
}

func NewViewEvent() scenariolib.JSONEvent {
	return scenariolib.JSONEvent{
		Type: "View",
		Arguments: map[string]interface{}{
			"offset": 0,
			"probability": 1,
			"docNo": -1,
		},
	}
}
