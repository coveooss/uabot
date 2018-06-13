package scenariolib

import (
	"errors"
	"math"
	"math/rand"

	"github.com/coveo/go-coveo/search"

	"encoding/json"
)

// ============== CLICK EVENT ======================
// =================================================

// ClickEvent a struct representing a click, it is defined by a clickRank, an
// offset and a probability to click.
type ClickEvent struct {
	Probability  float64                `json:"probability"`
	ClickRank    int                    `json:"docNo"`
	Offset       int                    `json:"offset,omitempty"`
	Quickview    bool                   `json:"quickview,omitempty"`
	CustomData   map[string]interface{} `json:"customData,omitempty"`
	FakeClick    bool                   `json:"fakeClick,omitempty"`
	FakeResponse json.RawMessage        `json:"fakeResponse,omitempty"`
}

// IsValid Validate a click event by applying different validation rules of dependant parameters etc.
func (click *ClickEvent) IsValid() (bool, string) {
	if click.Probability < 0 || click.Probability > 1 {
		return false, "A click event probability must be between 0 and 1."
	}

	if click.Offset < 0 {
		return false, "Offset must be a positive integer."
	}

	if click.ClickRank < -1 {
		return false, "Click rank must be -1 or >= 0."
	}

	if click.FakeClick && click.FakeResponse == nil {
		return false, "If you set parameter fakeClick to true, you must also send a fakeResponse."
	}
	return true, ""
}

// Execute Execute the click event, sending a click event to the usage analytics
func (click *ClickEvent) Execute(v *Visit) error {
	if click.FakeClick {
		searchUID := v.LastResponse.SearchUID
		fakeResponse := &search.Response{}
		if err := json.Unmarshal(click.FakeResponse, fakeResponse); err != nil {
			return err
		}
		v.LastResponse = fakeResponse
		v.LastResponse.SearchUID = searchUID
	}

	// Error handling, error if last response is nil, warning if last response had no results
	if v.LastResponse == nil {
		return errors.New("LastResponse is nil, execute a search first")
	}

	if click.ClickRank == -1 { // if rank == -1 we need to randomize a rank
		// Find a random rank within the possible click values accounting for the offset
		click.ClickRank = int(math.Abs(rand.NormFloat64()*2)) + click.Offset
	}

	// Make sure the click rank is not > the number of results.
	if v.LastResponse.TotalCount > 1 {
		topL := Min(v.LastQuery.NumberOfResults, v.LastResponse.TotalCount)
		click.ClickRank = Min(click.ClickRank, topL-1)
	} else {
		Warning.Printf("Last query %s returned no results cannot click", v.LastQuery.Q)
		return nil
	}

	if rand.Float64() <= click.Probability { // Probability to click
		if click.ClickRank > v.LastResponse.TotalCount {
			return errors.New("Click index out of bounds")
		}

		err := v.sendClickEvent(click.ClickRank, click.Quickview, click.CustomData)
		if err != nil {
			return err
		}
		return nil
	}
	Info.Printf("User chose not to click (probability %v%%)", int(click.Probability*100))
	return nil
}
