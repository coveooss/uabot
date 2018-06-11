package scenariolib_test

import (
	"encoding/json"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestPageViewEventValid(t *testing.T) {
	var testEventJson = []byte(`{"clickRank": 1, "probability": 0.5, "pageViewField": "pageViewFieldTest", "offset": 0, "contentType": "contentTypeTest", "customData": {"data1": "one"}}`)
	event := &scenariolib.ViewEvent{}

	// Test unmarshal json.
	if err := json.Unmarshal(testEventJson, event); err != nil {
		t.Errorf("Error Unmarshaling ViewEvent: %s", err)
	}

	if valid, message := event.IsValid(); !valid {
		t.Errorf("Expected ViewEvent.IsValid() to be true, got false (%s)", message)
	}

	if event.ClickRank != 1 {
		t.Errorf("Expected ViewEvent.ClickRank to be 1, got %d instead.", event.ClickRank)
	}

	if event.Probability != 0.5 {
		t.Errorf("Expected ViewEvent.Probability to be 2, got %f instead.", event.Probability)
	}

	if event.PageViewField != "pageViewFieldTest" {
		t.Errorf("Expected ViewEvent.PageViewField to be 'pageViewFieldTest', got %s instead.", event.PageViewField)
	}

	if event.ContentType != "contentTypeTest" {
		t.Errorf("Expected ViewEvent.ContentType to be 'contentTypeTest', got %s instead.", event.ContentType)
	}

	if event.CustomData == nil {
		t.Errorf("Expected ViewEvent.CustomData to not be nil, was nil instead.")
	}

	if event.CustomData["data1"] != "one" {
		t.Errorf("Expected ViewEvent.CustomData['data1'] to be 'one', was %s instead.", event.CustomData["data1"])
	}
}
