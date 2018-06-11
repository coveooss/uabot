// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import (
	"github.com/coveo/go-coveo/search"
)

// ============== FAKE SEARCH EVENT ======================
// =======================================================

// FakeSearchEvent a struct representing a search, is the response to set in the visit
type FakeSearchEvent struct {
	FakeResponse *search.Response `json:"fakeResponse"`
}

// IsValid Additional validation after the json unmarshal.
func (fakeSearch *FakeSearchEvent) IsValid() (bool, string) {
	return true, ""
}

// Execute the fake search event, set the Last response to the fake response
func (fakeSearch *FakeSearchEvent) Execute(v *Visit) error {
	v.LastQuery.Q = ""
	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}

	v.LastResponse = fakeSearch.FakeResponse
	v.LastResponse.SearchUID = resp.SearchUID
	return nil
}
