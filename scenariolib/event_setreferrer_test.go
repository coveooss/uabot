package scenariolib_test

import (
	"encoding/json"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestSetReferrerEvent(t *testing.T) {
	var testEventJson = []byte(`{"referrer": "testReferrer"}`)
	event := &scenariolib.SetReferrerEvent{}

	// Test unmarshal json.
	if err := json.Unmarshal(testEventJson, event); err != nil {
		t.Errorf("Error Unmarshaling SetReferrerEvent: %s", err)
	}

	if valid, message := event.IsValid(); !valid {
		t.Errorf("Expected SetReferrerEvent.IsValid() to be true, got false (%s)", message)
	}

	if event.Referrer != "testReferrer" {
		t.Errorf("Expected setReferrerEvent.Referrer to be 'testReferrer', got %s instead.", event.Referrer)
	}
}
