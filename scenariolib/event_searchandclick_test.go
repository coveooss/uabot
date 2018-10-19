package scenariolib_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestSearchAndClickEventValid(t *testing.T) {
	var testEventJson = []byte(`{"queryText": "queryTextTest", "probability": 0.5, "docClickTitle": "docTitleTest", "quickview": false, "caseSearch": false, "inputTitle": "inputTitleTest", "customData": {"data1": "one"}}`)
	event := &scenariolib.SearchAndClickEvent{}

	// Test unmarshal json.
	err := json.Unmarshal(testEventJson, event)
	ok(t, err)

	valid, message := event.IsValid()
	assert(t, valid, "Expected event to be valid, was false with error: %s", message)

	equals(t, "queryTextTest", event.Query)
	equals(t, 0.5, event.Probability)
	equals(t, "docTitleTest", event.DocTitle)
	assert(t, !event.Quickview, "Expected Quickview to be false.")
	assert(t, !event.CaseSearch, "Expected CaseSearch to be false.")
	equals(t, "inputTitleTest", event.InputTitle)

	// Expect CustomData to be not nil
	assert(t, event.CustomData != nil, "Expected CustomData to not be nil.")

	// Expect CustomData["data1"] to be "one"
	equals(t, "one", event.CustomData["data1"])
}

func TestDecorateSearchAndClickEvent(t *testing.T) {
	scenariolib.InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	expected := map[string]ExpectedRequest{
		"/rest/search/": {
			Method: "POST",
			Headers: map[string]string{
				"Authorization": "Bearer bot.searchToken",
				"Content-Type":  "application/json",
			},
			Body: []byte(`{
				"q": "aaaaaaaaaaa",
				"numberOfResults": 20,
				"tab": "All",
				"pipeline": "besttechCommunity"
			}`),
		},
		"/rest/v15/analytics/search/": {
			Method: "POST",
			Headers: map[string]string{
				"Authorization": "Bearer bot.analyticsToken",
				"Content-Type":  "application/json",
			},
			Body: []byte(`{
					"language": "en",
					"device": "Bot",
					"customData": {
						"JSUIVersion": "0.0.0.0;0.0.0.0",
						"customData1": "customValue 1",
						"ipaddress": "216.249.112.8"
					},
					"anonymous": true,
					"originLevel1": "BotSearch",
					"originLevel2": "",
					"searchQueryUid": "",
					"queryText": "aaaaaaaaaaa",
					"actionCause": "searchboxSubmit",
					"contextual": false
				}`),
		},
	}

	server := createMockServer(t, expected)
	defer server.Close() // Close the server when test finishes

	event := &scenariolib.SearchAndClickEvent{}
	conf, err := scenariolib.NewConfigFromPath("../scenarios_examples/TESTScenarios.json")
	ok(t, err)

	// Use the server url to define the endpoints
	conf.SearchEndpoint = server.URL + "/rest/search/"
	conf.AnalyticsEndpoint = server.URL + "/rest/v15/analytics/"

	v, err := scenariolib.NewVisit("bot.searchToken", "bot.analyticsToken", "scenario.UserAgent", "en", conf)

	v.SetupGeneral()
	event.Execute(v)
}
