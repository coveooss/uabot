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
	err := json.Unmarshal(testEventJson, event)
	ok(t, err)

	valid, message := event.IsValid()
	assert(t, valid, "Expected event to be valid, was false with error: %s", message)

	equals(t, "nameTest", event.Name)

	equals(t, "@uri", event.ConstantExpression)
}
