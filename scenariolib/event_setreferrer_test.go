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
	err := json.Unmarshal(testEventJson, event)
	ok(t, err)

	valid, message := event.IsValid()
	assert(t, valid, "Expected event to be valid, was false with error: %s", message)

	equals(t, "testReferrer", event.Referrer)
}
