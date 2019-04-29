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
	err := json.Unmarshal(testEventJson, event)
	ok(t, err)

	valid, message := event.IsValid()
	assert(t, valid, "Expected CustomEvent.IsValid() to be true, was false with error: %s", message)

	equals(t, "eventTypeTest", event.EventType)

	equals(t, "eventValueTest", event.EventValue)

	// Expect CustomData to be not nil
	assert(t, event.CustomData != nil, "Expected CustomData to not be nil.")

	// Expect CustomData["data1"] to be "one"
	equals(t, "one", event.CustomData["data1"])
}
