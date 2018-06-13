package scenariolib_test

import (
	"encoding/json"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestSearchAndClickEventValid(t *testing.T) {
	var testEventJson = []byte(`{"query": "queryTextTest", "probability": 0.5, "docClickTitle": "docTitleTest", "quickview": false, "caseSearch": false, "inputTitle": "inputTitleTest", "customData": {"data1": "one"}}`)
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
