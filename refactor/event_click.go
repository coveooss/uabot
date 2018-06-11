package refactor

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"

	"github.com/go-coveo/analytics"
	"github.com/go-coveo/search"
)

// ============== CLICK EVENT ======================
// ==================================================

// A ClickEvent is an event sent when the user clicks on a document
type ClickEvent struct {
	DocNo            int                    `json:"docNo,omitempty"`
	Offset           int                    `json:"offset,omitempty"`
	Probability      float64                `json:"probability"`
	Quickview        bool                   `json:"quickview,omitempty"`
	CustomData       map[string]interface{} `json:"customData,omitempty"`
	FakeClick        bool                   `json:"fakeClick,omitempty"`
	FakeResponseJSON json.RawMessage        `json:"fakeResponse,omitempty"`
	fakeResponse     *search.Response
}

// Parse the remaining bits of the json event into the right arguments for this event.
func (e *ClickEvent) Parse(jse *JSONEvent) error {
	if err := json.Unmarshal(jse.Arguments, e); err != nil {
		return err
	}
	if e.Probability < 0 || e.Probability > 1 {
		return errors.New("Probability must be between 0 and 1 in a click event")
	}
	if e.Offset < 0 {
		return errors.New("Offset must be positive in a click event")
	}
	if e.DocNo < -1 {
		return errors.New("DocNo must be > 0 or -1 (for a random rank) in a click event")
	}
	if e.FakeClick {
		if err := json.Unmarshal(e.FakeResponseJSON, e.fakeResponse); err != nil {
			return errors.New("FakeResponse must be a search.Response.")
		}
	}
	return nil
}

// Execute Execute the search event, runs the query and sends a search event to
// the analytics.
func (e *ClickEvent) Execute(v *Visit) error {
	if e.FakeClick {
		v.LastResponse = e.fakeResponse
	}
	if v.LastResponse == nil {
		return errors.New("LastResponse is nil cannot execute a click event.")
	}
	if v.LastResponse.TotalCount < 1 {
		// Warning.Printf("Last query %s returned no results cannot click", v.LastQuery.Q)
		return nil
	}

	if rand.Float64() <= e.Probability {
		if e.DocNo == -1 {
			e.DocNo = e.findClickRank(v)
		}

		if err := e.sendClickEvent(v); err != nil {
			return err
		}
	}

	// Info.Printf("User chose not to click (probability %v%%)", int(ce.probability*100))
	return nil
}

func (e *ClickEvent) sendClickEvent(v *Visit) error {
	if v.LastResponse == nil {
		return errors.New("LastResponse was nil cannot send click event.")
	}
	// Info.Printf("Sending ClickEvent rank=%d (quickview %v)", rank+1, quickview)
	var validcast bool
	result := v.LastResponse.Results[e.DocNo]

	event := analytics.NewClickEvent()

	if event.DocumentURIHash, validcast = result.Raw["sysurihash"].(string); !validcast {
		return errors.New("Result.raw.sysurihash is not a string")
	}
	if event.CollectionName, validcast = result.Raw["syscollection"].(string); !validcast {
		return errors.New("Result.raw.syscollection is not a string")
	}
	if event.SourceName, validcast = result.Raw["syssource"].(string); !validcast {
		return errors.New("Result.raw.syssource is not a string")
	}
	if e.Quickview {
		event.ActionCause = "documentQuickview"
		event.ViewMethod = "documentQuickview"
	} else {
		event.ActionCause = "documentOpen"
	}
	event.DocumentURI = result.URI
	event.SearchQueryUID = v.LastResponse.SearchUID
	event.DocumentPosition = e.DocNo + 1
	event.DocumentTitle = result.Title
	event.QueryPipeline = v.LastResponse.Pipeline
	event.DocumentURL = result.ClickURI
	event.Username = v.User.Email
	event.Anonymous = v.User.Anonymous
	event.OriginLevel1 = v.OriginLevel1
	event.OriginLevel2 = v.OriginLevel2
	event.Language = v.User.Language

	// CustomData
	defaultCustomData := map[string]interface{}{
		"JSUIVersion": JSUIVERSION,
		"ipadress":    v.User.IP,
		"author":      v.Config.RandomDocumentAuthors[rand.Intn(len(v.Config.RandomDocumentAuthors))],
	}
	if event.CustomData == nil {
		event.CustomData = defaultCustomData
	} else {
		for k, v := range defaultCustomData {
			if event.CustomData[k] == nil {
				event.CustomData[k] = v
			}
		}
	}

	// Send all the possible random custom data that can be added from the config
	// scenario file.
	for _, elem := range v.Config.RandomCustomData {
		event.CustomData[elem.APIName] = elem.Values[rand.Intn(len(elem.Values))]
	}

	//Trace.Printf("Sending ClickEvent [ rank=%d quickview=%v ]", e.ClickRank+1, e.Quickview)
	if err := v.UAClient.SendClickEvent(event); err != nil {
		return err
	}

	return nil
}

func (e *ClickEvent) findClickRank(v *Visit) int {
	var clickRank int
	if v.LastResponse.TotalCount > 1 {
		topL := Min(v.LastQuery.NumberOfResults, v.LastResponse.TotalCount)
		rndRank := int(math.Abs(rand.NormFloat64()*2)) + e.Offset
		clickRank = Min(rndRank, topL-1)
	}
	return Min(clickRank, v.LastResponse.TotalCount)
}

// Check for interface implementation
var _ Event = (*ClickEvent)(nil)
