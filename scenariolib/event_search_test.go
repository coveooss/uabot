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
	if err := json.Unmarshal(testEventJson, event); err != nil {
		t.Errorf("Error Unmarshaling SearchEvent: %s", err)
	}

	if valid, message := event.IsValid(); !valid {
		t.Errorf("Expected SearchEvent.IsValid() to be true, got false (%s)", message)
	}

	if event.Query != "queryTextTest" {
		t.Errorf("Expected SearchEvent.Query to be 'queryTextTest', got %s instead.", event.Query)
	}

	if event.IgnoreEvent != false {
		t.Errorf("Expected SearchEvent.IgnoreEvent to be false, got %t instead.", event.IgnoreEvent)
	}

	if event.GoodQuery != true {
		t.Errorf("Expected SearchEvent.GoodQuery to be true, got %t instead.", event.GoodQuery)
	}

	if event.ActionCause != "actionCauseTest" {
		t.Errorf("Expected SearchEvent.ActionCause to be 'actionCauseTest', got %s instead.", event.ActionCause)
	}

	if event.CaseSearch != false {
		t.Errorf("Expected SearchEvent.CaseSearch to be false, got %t instead.", event.CaseSearch)
	}

	if event.InputTitle != "inputTitleTest" {
		t.Errorf("Expected SearchEvent.InputTitle to be 'inputTitleTest', got %s instead.", event.InputTitle)
	}

	if event.MatchLanguage != false {
		t.Errorf("Expected SearchEvent.MatchLanguage to be false, got %t instead.", event.MatchLanguage)
	}

	if event.CustomData == nil {
		t.Errorf("Expected SearchEvent.CustomData to not be nil, was nil instead.")
	}

	if event.CustomData["data1"] != "one" {
		t.Errorf("Expected SearchEvent.CustomData['data1'] to be 'one', was %s instead.", event.CustomData["data1"])
	}
}
