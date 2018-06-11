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
	if err := json.Unmarshal(testEventJson, event); err != nil {
		t.Errorf("Error Unmarshaling FacetEvent: %s", err)
	}

	if valid, message := event.IsValid(); !valid {
		t.Errorf("Expected FacetEvent.IsValid() to be true, got false (%s)", message)
	}

	if event.FacetTitle != "facetTitleTest" {
		t.Errorf("Expected FacetEvent.FacetTitle to be 'facetTitleTest', got %s instead.", event.FacetTitle)
	}

	if event.FacetValue != "facetValueTest" {
		t.Errorf("Expected FacetEvent.FacetValue to be 'facetValueTest', got %s instead.", event.FacetValue)
	}

	if event.FacetField != "facetFieldTest" {
		t.Errorf("Expected FacetEvent.FacetField to be 'facetFieldTest', got %s instead.", event.FacetField)
	}

	if event.CustomData == nil {
		t.Errorf("Expected FacetEvent.CustomData to not be nil, was nil instead.")
	}

	if event.CustomData["data1"] != "one" {
		t.Errorf("Expected FacetEvent.CustomData['data1'] to be 'one', was %s instead.", event.CustomData["data1"])
	}
}
