package scenariolib_test

import (
	"encoding/json"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestTabChangeEvent(t *testing.T) {
	var testEventJson = []byte(`{"name": "nameTest", "cq": "@uri"}`)
	event := &scenariolib.TabChangeEvent{}

	// Test unmarshal json.
	if err := json.Unmarshal(testEventJson, event); err != nil {
		t.Errorf("Error Unmarshaling SetReferrerEvent: %s", err)
	}

	if valid, message := event.IsValid(); !valid {
		t.Errorf("Expected SetReferrerEvent.IsValid() to be true, got false (%s)", message)
	}

	if event.Name != "nameTest" {
		t.Errorf("Expected TabChangeEvent.Name to be 'nameTest', got %s instead.", event.Name)
	}

	if event.ConstantExpression != "@uri" {
		t.Errorf("Expected TabChangeEvent.ConstantExpression to be '@uri', got %s instead.", event.ConstantExpression)
	}
}
