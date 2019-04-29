package scenariolib_test

import (
	"encoding/json"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestSearchEventValid(t *testing.T) {
	var testEventJson = []byte(`{"queryText": "queryTextTest", "ignoreEvent": false, "goodQuery": true, "actionCause": "actionCauseTest", "caseSearch": false, "inputTitle": "inputTitleTest", "matchLanguage": false, "customData": {"data1": "one"}}`)
	event := &scenariolib.SearchEvent{}

	// Test unmarshal json.
	err := json.Unmarshal(testEventJson, event)
	ok(t, err)

	valid, message := event.IsValid()
	assert(t, valid, "Expected event to be valid, was false with error: %s", message)

	equals(t, "queryTextTest", event.Query)

	assert(t, !event.IgnoreEvent, "Expected IgnoreEvent to be false.")

	assert(t, event.GoodQuery, "Expected GoodQuery to be true.")

	equals(t, "actionCauseTest", event.ActionCause)

	assert(t, !event.CaseSearch, "Expected CaseSearch to be false.")

	equals(t, "inputTitleTest", event.InputTitle)

	assert(t, !event.MatchLanguage, "Expected MatchLanguage to be false.")

	// Expect CustomData to be not nil
	assert(t, event.CustomData != nil, "Expected CustomData to not be nil.")

	// Expect CustomData["data1"] to be "one"
	equals(t, "one", event.CustomData["data1"])
}
