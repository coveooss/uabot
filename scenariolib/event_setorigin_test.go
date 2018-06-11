package scenariolib_test

import (
	"encoding/json"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestSetOriginEvent(t *testing.T) {
	var testEventJson = []byte(`{"originLevel1": "testOriginLevel1", "originLevel2": "testOriginLevel2", "originLevel3": "testOriginLevel3"}`)
	event := &scenariolib.SetOriginEvent{}

	// Test unmarshal json.
	if err := json.Unmarshal(testEventJson, event); err != nil {
		t.Errorf("Error Unmarshaling SetOriginEvent: %s", err)
	}

	if valid, message := event.IsValid(); !valid {
		t.Errorf("Expected SetOriginEvent.IsValid() to be true, got false (%s)", message)
	}

	if event.OriginLevel1 != "testOriginLevel1" {
		t.Errorf("Expected setOriginEvent.OriginLevel1 to be 'testOriginLevel1', got %s instead.", event.OriginLevel1)
	}

	if event.OriginLevel2 != "testOriginLevel2" {
		t.Errorf("Expected setOriginEvent.OriginLevel2 to be 'testOriginLevel2', got %s instead.", event.OriginLevel2)
	}

	if event.OriginLevel3 != "testOriginLevel3" {
		t.Errorf("Expected setOriginEvent.OriginLevel3 to be 'testOriginLevel3', got %s instead.", event.OriginLevel3)
	}
}
