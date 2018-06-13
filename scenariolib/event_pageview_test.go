package scenariolib_test

import (
	"encoding/json"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestPageViewEventValid(t *testing.T) {
	var testEventJson = []byte(`{"clickRank": 1, "probability": 0.5, "pageViewField": "pageViewFieldTest", "offset": 0, "contentType": "contentTypeTest", "customData": {"data1": "one"}}`)
	event := &scenariolib.ViewEvent{}

	// Test unmarshal json.
	err := json.Unmarshal(testEventJson, event)
	ok(t, err)

	valid, message := event.IsValid()
	assert(t, valid, "Expected event to be valid, was false with error: %s", message)

	equals(t, 1, event.ClickRank)

	equals(t, 0.5, event.Probability)

	equals(t, "pageViewFieldTest", event.PageViewField)

	equals(t, "contentTypeTest", event.ContentType)

	// Expect CustomData to be not nil
	assert(t, event.CustomData != nil, "Expected CustomData to not be nil.")

	// Expect CustomData["data1"] to be "one"
	equals(t, "one", event.CustomData["data1"])
}
