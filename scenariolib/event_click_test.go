package scenariolib_test

import (
	"encoding/json"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestClickEventValid(t *testing.T) {
	var testEventJson = []byte(`{"probability": 1, "docNo": -1, "offset": 2, "quickview": false, "customData": {"data1": "one"}}`)
	event := &scenariolib.ClickEvent{}

	// Test unmarshal json.
	if err := json.Unmarshal(testEventJson, event); err != nil {
		t.Errorf("Error Unmarshaling ClickEvent: %s", err)
	}

	if valid, message := event.IsValid(); !valid {
		t.Errorf("Expected ClickEvent.IsValid() to be true, got false (%s)", message)
	}

	if event.Probability != 1 {
		t.Errorf("Expected clickEvent.Probability to be 1, got %f instead.", event.Probability)
	}

	if event.ClickRank != -1 {
		t.Errorf("Expected clickEvent.ClickRank to be -1, got %d instead.", event.ClickRank)
	}

	if event.Offset != 2 {
		t.Errorf("Expected clickEvent.Offset to be 2, got %d instead.", event.Offset)
	}

	if event.Quickview != false {
		t.Errorf("Expected clickEvent.Quickview to be false, got %t instead.", event.Quickview)
	}

	if event.CustomData == nil {
		t.Errorf("Expected clickEvent.CustomData to be a map[string]interface{}, was nil instead.")
	}

	if event.CustomData["data1"] != "one" {
		t.Errorf("Expected clickEvent.CustomData[data1] to be 'one', got %v instead.", event.CustomData["data1"])
	}
}
