package scenariolib_test

import (
	"encoding/json"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestSearchAndClickEventValid(t *testing.T) {
	var testEventJson = []byte(`{"query": "queryTextTest", "probability": 0.5, "docTitle": "docTitleTest", "quickview": false, "caseSearch": false, "inputTitle": "inputTitleTest", "customData": {"data1": "one"}}`)
	event := &scenariolib.SearchAndClickEvent{}

	// Test unmarshal json.
	if err := json.Unmarshal(testEventJson, event); err != nil {
		t.Errorf("Error Unmarshaling SearchAndClickEvent: %s", err)
	}

	if valid, message := event.IsValid(); !valid {
		t.Errorf("Expected SearchAndClickEvent.IsValid() to be true, got false (%s)", message)
	}

	if event.Query != "queryTextTest" {
		t.Errorf("Expected SearchAndClickEvent.Query to be 'queryTextTest', got %s instead.", event.Query)
	}

	if event.Probability != 0.5 {
		t.Errorf("Expected SearchAndClickEvent.Probability to be 0.5, got %f instead.", event.Probability)
	}

	if event.DocTitle != "docTitleTest" {
		t.Errorf("Expected SearchAndClickEvent.DocTitle to be 'docTitleTest', got %s instead.", event.DocTitle)
	}

	if event.Quickview != false {
		t.Errorf("Expected SearchAndClickEvent.Quickview to be false, got %t instead.", event.Quickview)
	}

	if event.CaseSearch != false {
		t.Errorf("Expected SearchAndClickEvent.CaseSearch to be false, got %t instead.", event.CaseSearch)
	}

	if event.InputTitle != "inputTitleTest" {
		t.Errorf("Expected SearchAndClickEvent.InputTitle to be 'inputTitleTest', got %s instead.", event.InputTitle)
	}

	if event.CustomData == nil {
		t.Errorf("Expected SearchEvent.CustomData to not be nil, was nil instead.")
	}

	if event.CustomData["data1"] != "one" {
		t.Errorf("Expected SearchEvent.CustomData['data1'] to be 'one', was %s instead.", event.CustomData["data1"])
	}
}
