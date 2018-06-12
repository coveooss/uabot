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
	err := json.Unmarshal(testEventJson, event)
	ok(t, err)

	valid, message := event.IsValid()
	assert(t, valid, "Expected event to be valid, was false with error: %s", message)

	equals(t, "testOriginLevel1", event.OriginLevel1)

	equals(t, "testOriginLevel2", event.OriginLevel2)

	equals(t, "testOriginLevel3", event.OriginLevel3)
}
