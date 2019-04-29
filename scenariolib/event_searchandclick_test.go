package scenariolib_test

import (
	"encoding/json"
	"os"
	"testing"

	"../defaults"
	"../scenariolib"
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

	// All requests caught by the Test server will be added to 'requests'
	requests := make(map[string]RestRequest)

	server := createTestServer(t, requests)
	defer server.Close() // Close the server when test finishes

	event := &scenariolib.SearchAndClickEvent{}
	conf, err := scenariolib.NewConfigFromPath("../scenarios_examples/TESTScenarios.json")
	ok(t, err)

	// Use the server url to define the endpoints
	conf.SearchEndpoint = server.URL + defaults.SEARCH_REST_PATH
	conf.AnalyticsEndpoint = server.URL + defaults.ANALYTICS_REST_PATH

	v, err := scenariolib.NewVisit("bot.searchToken", "bot.analyticsToken", "scenario.UserAgent", "en", conf)

	v.SetupGeneral()
	event.Execute(v)

	// Validate the search request we expect Execute() to send.
	req, exists := requests["/rest/search/"]
	assert(t, exists, "Missing request for /rest/search/")
	equals(t, "POST", req.Method)
	equals(t, "Bearer bot.searchToken", req.Headers["Authorization"])
	equals(t, "application/json", req.Headers["Content-Type"])
	expectedBody := []byte(`{
		"q": "aaaaaaaaaaa",
		"numberOfResults": 20,
		"tab": "All",
		"pipeline": "besttechCommunity"
	}`)
	eq, err := JSONBytesEqual(expectedBody, req.Body)
	assert(t, eq, "The Request's body for Search is not what we expected\nGot: %s\nExp: %s", expectedBody, req.Body)

	// Validate analytics request we expect Execute() to send.
	req, exists = requests["/rest/v15/analytics/search/"]
	assert(t, exists, "Missing request for /rest/v15/analytics/search/")
	equals(t, "POST", req.Method)
	equals(t, "Bearer bot.analyticsToken", req.Headers["Authorization"])
	equals(t, "application/json", req.Headers["Content-Type"])
	expectedBody = []byte(`{
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
	}`)
	eq, err = JSONBytesEqual(expectedBody, req.Body)
	assert(t, eq, "The Request's body for Analytics is not what we expected\nGot: %s\nExp: %s", expectedBody, req.Body)
}

func TestDecorateSearchAndClickEvent2(t *testing.T) {
	scenariolib.InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	// All requests caught by the Test server will be added to 'requests'
	requests := make(map[string]RestRequest)

	server := createTestServer(t, requests)
	defer server.Close() // Close the server when test finishes

	event := &scenariolib.SearchAndClickEvent{}
	conf, err := scenariolib.NewConfigFromPath("../scenarios_examples/DemoMovies.json")
	ok(t, err)

	// Use the server url to define the endpoints
	conf.SearchEndpoint = server.URL + defaults.SEARCH_REST_PATH
	conf.AnalyticsEndpoint = server.URL + defaults.ANALYTICS_REST_PATH

	v, err := scenariolib.NewVisit("bot.searchToken", "bot.analyticsToken", "scenario.UserAgent", "en", conf)

	v.SetupGeneral()
	event.Execute(v)

	// Validate the search request we expect Execute() to send.
	req, exists := requests["/rest/search/"]
	assert(t, exists, "Missing request for /rest/search/")
	equals(t, "POST", req.Method)
	equals(t, "Bearer bot.searchToken", req.Headers["Authorization"])
	equals(t, "application/json", req.Headers["Content-Type"])
	expectedBody := []byte(`{
		"q": "Gostbuster",
		"numberOfResults": 20,
		"tab": "All",
		"pipeline": "ML"
	}`)
	eq, err := JSONBytesEqual(expectedBody, req.Body)
	assert(t, eq, "The Request's body for Search is not what we expected\nGot: %s\nExp: %s", expectedBody, req.Body)

	// Validate analytics request we expect Execute() to send.
	req, exists = requests["/rest/v15/analytics/search/"]
	assert(t, exists, "Missing request for /rest/v15/analytics/search/")
	equals(t, "POST", req.Method)
	equals(t, "Bearer bot.analyticsToken", req.Headers["Authorization"])
	equals(t, "application/json", req.Headers["Content-Type"])
	expectedBody = []byte(`{
		"language": "en",
		"device": "Bot",
		"customData": {
			"JSUIVersion": "0.0.0.0;0.0.0.0",
			"c_isbot": "true",
			"ipaddress": "198.199.154.209"
		},
		"username": "avery.caldwell@hexzone.com",
		"originLevel1": "Movie",
		"originLevel2": "default",
		"searchQueryUid": "",
		"queryText": "Gostbuster",
		"actionCause": "searchboxSubmit",
		"contextual": false
	}`)
	eq, err = JSONBytesEqual(expectedBody, req.Body)
	assert(t, eq, "The Request's body for Analytics is not what we expected\nGot: %s\nExp: %s", expectedBody, req.Body)
}
