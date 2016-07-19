// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import (
	"encoding/json"
	"errors"
	"github.com/coveo/go-coveo/search"
)

// ============== FAKE SEARCH EVENT ======================
// ==================================================

// FakeSearchEvent a struct representing a search, is the response to set in the visit
type FakeSearchEvent struct {
	fakeResponse *search.Response
}

func newFakeSearchEvent(e *JSONEvent, c *Config) (*FakeSearchEvent, error) {
	se := new(FakeSearchEvent)

	if e.Arguments["fakeResponse"] != nil {
		jsonFalseResponse, _ := json.Marshal(e.Arguments["fakeResponse"])
		err := json.Unmarshal(jsonFalseResponse, &se.fakeResponse)
		if err != nil {
			return nil, errors.New("Parameter falseResponse must be a search.Response")
		}
	}

	return se, nil
}

// Execute Execute the fake search event, set the Last response to the fake response
func (fse *FakeSearchEvent) Execute(v *Visit) error {
	v.LastQuery.Q = ""
	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}

	v.LastResponse = fse.fakeResponse
	v.LastResponse.SearchUID = resp.SearchUID
	return nil
}
