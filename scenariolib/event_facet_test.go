package scenariolib_test

import (
	"encoding/json"
	"testing"

	"github.com/coveo/uabot/scenariolib"
)

func TestFacetEventValid(t *testing.T) {
	var testEventJson = []byte(`{"facetTitle": "facetTitleTest", "facetValue": "facetValueTest", "facetField": "facetFieldTest", "customData": {"data1": "one"}}`)
	event := &scenariolib.FacetEvent{}

	// Test unmarshal json.
	err := json.Unmarshal(testEventJson, event)
	ok(t, err)

	valid, message := event.IsValid()
	assert(t, valid, "Expected event to be valid, was false with error: %s", message)

	equals(t, "facetTitleTest", event.FacetTitle)

	equals(t, "facetValueTest", event.FacetValue)

	equals(t, "facetFieldTest", event.FacetField)

	// Expect CustomData to be not nil
	assert(t, event.CustomData != nil, "Expected CustomData to not be nil.")

	// Expect CustomData["data1"] to be "one"
	equals(t, "one", event.CustomData["data1"])
}
