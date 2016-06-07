package visit

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"

	"github.com/coveo/go-coveo/analytics"
)

// Click Contains the arguments necessary to send a click event to the analytics
type Click struct {
	ClickRank   int     `json:"docNo"`
	Offset      int     `json:"offset"`
	Probability float64 `json:"probability"`
	Quickview   bool    `json:"quickview,omitempty"`
}

// Parse Parse the different arguments in the JSONEvent to build the click event
func (e *Click) Parse(jse *JSONEvent) error {
	err := json.Unmarshal(jse.Arguments, e)
	if err != nil {
		return err
	}

	if e.Offset < 0 {
		return errors.New("Offset cannot be negative")
	}

	if e.Probability < 0 || e.Probability > 1 {
		return errors.New("Probability is out of bounds")
	}
	return nil
}

// Execute Send the event to the analytics endpoint
func (e *Click) Execute(v *Visit) error {
	if v.LastResponse.TotalCount < 1 {
		Warning.Printf("Last query %s returned no results cannot click", v.LastQuery.Q)
		return nil
	}

	roll := rand.Float64()
	if roll <= e.Probability {
		if e.ClickRank == -1 {
			e.ClickRank = 0
			// Find a random rank within the possible click values accounting for the offset
			if v.LastResponse.TotalCount > 1 {
				topL := Min(v.LastQuery.NumberOfResults, v.LastResponse.TotalCount)
				rndRank := int(math.Abs(rand.NormFloat64()*2)) + e.Offset
				e.ClickRank = Min(rndRank, topL-1)
			}
		}
		if e.ClickRank > v.LastResponse.TotalCount {
			e.ClickRank = v.LastResponse.TotalCount
		}

		if err := e.sendClickEvent(v); err != nil {
			return err
		}
		return nil
	}
	Info.Printf("User chose not to click [ roll=%d%% , probability=%d%% ]", int(roll*100), int(e.Probability*100))
	return nil
}

func (e *Click) sendClickEvent(v *Visit) error {
	event, err := analytics.NewClickEvent()
	if err != nil {
		return err
	}
	var validCast bool

	if event.DocumentURIHash, validCast = v.LastResponse.Results[e.ClickRank].Raw["sysurihash"].(string); !validCast {
		return errors.New("Result.raw.sysurihash is not a string")
	}
	if event.CollectionName, validCast = v.LastResponse.Results[e.ClickRank].Raw["syscollection"].(string); !validCast {
		return errors.New("Result.raw.syscollection is not a string")
	}
	if event.SourceName, validCast = v.LastResponse.Results[e.ClickRank].Raw["syssource"].(string); !validCast {
		return errors.New("Result.raw.syssource is not a string")
	}
	if e.Quickview {
		event.ActionCause = "documentQuickview"
		event.ViewMethod = "documentQuickview"
	} else {
		event.ActionCause = "documentOpen"
	}
	event.DocumentURI = v.LastResponse.Results[e.ClickRank].URI
	event.SearchQueryUID = v.LastResponse.SearchUID
	event.DocumentPosition = e.ClickRank + 1
	event.DocumentTitle = v.LastResponse.Results[e.ClickRank].Title
	event.QueryPipeline = v.LastResponse.Pipeline
	event.DocumentURL = v.LastResponse.Results[e.ClickRank].ClickUri
	event.Username = v.User.Email
	event.OriginLevel1 = v.OriginLevel1
	event.OriginLevel2 = v.OriginLevel2
	event.Anonymous = v.User.Anonymous
	event.Language = v.User.Language

	event.CustomData = make(map[string]interface{})
	event.CustomData["JSUIVersion"] = JSUIVERSION
	event.CustomData["ipadress"] = v.User.IP
	event.CustomData["author"] = v.Config.RandomDocumentAuthors[rand.Intn(len(v.Config.RandomDocumentAuthors))]

	if v.Config.AllowEntitlements { // Custom shit for besttech, I don't like it
		event.CustomData["entitlement"] = generateEntitlementBesttech(v.User.Anonymous)
	}

	// Send all the possible random custom data that can be added from the config
	// scenario file.
	for _, elem := range v.Config.RandomCustomData {
		event.CustomData[elem.APIName] = elem.Values[rand.Intn(len(elem.Values))]
	}

	Trace.Printf("Sending ClickEvent [ rank=%d quickview=%v ]", e.ClickRank+1, e.Quickview)
	err = v.UAClient.SendClickEvent(event)
	if err != nil {
		return err
	}
	return nil
}

func generateEntitlementBesttech(isAnonymous bool) string {
	if isAnonymous {
		return "Anonymous"
	}
	if rand.Float64() <= 0.1 {
		return "Premier"
	}
	return "Basic"
}

// Check for interface implementation
var _ Event = (*Click)(nil)
