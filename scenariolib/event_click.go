package scenariolib

import (
	"errors"
	"math"
	"math/rand"

	"github.com/coveooss/go-coveo/search"

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
	if v.LastResponse.TotalCount < 1 {
		Warning.Printf("Last query %s returned no results cannot send view event", v.LastQuery.Q)
		return nil
	}

	if rand.Float64() <= click.Probability { // Probability to click
		click.ClickRank = computeClickRank(v, click.ClickRank, click.Offset)

		// We leave this option because it means voluntarily someone set a clickRank > number of results for his query.
		if click.ClickRank > v.LastResponse.TotalCount {
			Warning.Printf("PageView index out of bounds, not sending event")
			return nil
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

// Randomize a click rank if the clickRank is -1
func computeClickRank(v *Visit, clickRank, offset int) (computedRank int) {
	computedRank = clickRank
	if computedRank == -1 { // if rank == -1 we need to randomize a rank
		// Default to the first result to click on
		computedRank = 0
		// Find a random rank within the possible click values accounting for the offset
		if v.LastResponse.TotalCount > 1 {
			topL := Min(v.LastQuery.NumberOfResults, v.LastResponse.TotalCount)
			rndRank := int(math.Abs(rand.NormFloat64()*2)) + offset
			computedRank = Min(rndRank, topL-1)
		}
	}
	return
}
