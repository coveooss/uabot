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
	err := json.Unmarshal(testEventJson, event)
	ok(t, err)

	// Test if the event is valid.
	valid, message := event.IsValid()
	assert(t, valid, "Expected ClickEvent.IsValid() to be true, was false with error: %s", message)

	// Expect Probability to be 1.0
	equals(t, 1.0, event.Probability)

	// Expect ClickRank to be -1
	equals(t, -1, event.ClickRank)

	// Expect Offset to be 2
	equals(t, 2, event.Offset)

	// Expect Quickview to be false
	assert(t, !event.Quickview, "Expected Quickview to be false")

	// Expect CustomData to be not nil
	assert(t, event.CustomData != nil, "Expected CustomData to not be nil.")

	// Expect CustomData["data1"] to be "one"
	equals(t, "one", event.CustomData["data1"])
}
