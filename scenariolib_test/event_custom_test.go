package scenariolib_test

import (
	"encoding/json"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestCustomEventValid(t *testing.T) {
	var testEventJson = []byte(`{"eventType": "eventTypeTest", "eventValue": "eventValueTest", "customData": {"data1": "one"}}`)
	event := &scenariolib.CustomEvent{}

	// Test unmarshal json.
	if err := json.Unmarshal(testEventJson, event); err != nil {
		t.Errorf("Error Unmarshaling CustomEvent: %s", err)
	}

	if valid, message := event.IsValid(); !valid {
		t.Errorf("Expected CustomEvent.IsValid() to be true, got false (%s)", message)
	}

	if event.EventType != "eventTypeTest" {
		t.Errorf("Expected CustomEvent.EventType to be 'eventTypeTest', got %s instead.", event.EventType)
	}

	if event.EventValue != "eventValueTest" {
		t.Errorf("Expected CustomEvent.EventValue to be 'eventValueTest', got %s instead.", event.EventValue)
	}

	if event.CustomData == nil {
		t.Errorf("Expected CustomEvent.CustomData to not be nil, was nil instead.")
	}

	if event.CustomData["data1"] != "one" {
		t.Errorf("Expected CustomEvent.CustomData['data1'] to be 'one', was %s instead.", event.CustomData["data1"])
	}
}
