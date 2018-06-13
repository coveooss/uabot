package scenariolib

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
)

// ============== SEARCH AND CLICK EVENT ======================
// ============================================================

// SearchAndClickEvent represents a search event followed by a click on a specific
// document found by the title
type SearchAndClickEvent struct {
	Query        string                 `json:"query"`
	Probability  float64                `json:"probability"`
	DocTitle     string                 `json:"docClickTitle,omitempty"`
	MatchField   string                 `json:"matchField,,omitempty"`
	MatchPattern string                 `json:"matchPattern,omitempty"`
	Quickview    bool                   `json:"quickview,omitempty"`
	CaseSearch   bool                   `json:"caseSearch,omitempty"`
	InputTitle   string                 `json:"inputTitle,omitempty"`
	CustomData   map[string]interface{} `json:"customData,omitempty"`
	RegexMatch   *regexp.Regexp
}

// IsValid Additional validation after the json unmarshal. And compilation of the regex if available.
func (searchClick *SearchAndClickEvent) IsValid() (bool, string) {
	if searchClick.DocTitle == "" {
		if searchClick.MatchField == "" || searchClick.MatchPattern == "" {
			return false, "If you are not using a [docClickTitle] you must provide both [matchField and matchPattern]"
		}
	} else {
		if searchClick.MatchField != "" || searchClick.MatchPattern != "" {
			return false, "If you provide a [docClickTitle] you cannot also use [matchField and/or matchPattern]"
		}
	}
	var err error
	if searchClick.RegexMatch, err = regexp.Compile(searchClick.MatchPattern); err != nil {
		return false, "Failed to compile regex pattern : " + err.Error()
	}

	return true, ""
}

// Execute the search and click event sending both events to the analytics
func (searchClick *SearchAndClickEvent) Execute(v *Visit) error {
	Info.Printf("Executing SearchAndClickEvent : Searching for %s, clicking on %s (quickview %v)", searchClick.Query, searchClick.DocTitle, searchClick.Quickview)
	// Execute the search event
	search := new(SearchEvent)
	search.Query = searchClick.Query
	search.Keyword = searchClick.Query
	search.CustomData = make(map[string]interface{})
	if searchClick.CaseSearch {
		search.Query = fmt.Sprintf("($some(keywords: %s, match: 1, removeStopWords: true, maximum: 300)) ($sort(criteria: relevancy))", search.Query)
		search.ActionCause = "inputChange"
		search.ActionType = "caseCreation"
		search.CustomData["inputTitle"] = searchClick.InputTitle
	} else {
		search.ActionCause = "searchboxSubmit"
		search.ActionType = "search box"
	}
	// Override possible values of customData with the specific customData sent
	for k, v := range searchClick.CustomData {
		search.CustomData[k] = v
	}
	if err := search.Execute(v); err != nil {
		return err
	}

	if v.LastResponse.TotalCount < 1 {
		return errors.New("Last query returned no results")
	}

	var timeToWait int
	if v.Config.TimeBetweenActions > 0 {
		timeToWait = v.Config.TimeBetweenActions
	} else {
		timeToWait = DEFAULTTIMEBETWEENACTIONS
	}
	WaitBetweenActions(timeToWait, v.Config.IsWaitConstant)

	if rand.Float64() <= searchClick.Probability {
		var rank int
		if searchClick.MatchField != "" {
			rank = v.FindDocumentRankByMatchingField(searchClick.MatchField, searchClick.RegexMatch)
		} else {
			rank = v.FindDocumentRankByTitle(searchClick.DocTitle)
		}
		if rank >= 0 {
			Info.Printf("Sending ClickEvent => Found document at rank : %d", rank+1)

			click := new(ClickEvent)
			click.ClickRank = rank
			click.Offset = 0
			click.Probability = 1
			click.Quickview = searchClick.Quickview

			click.CustomData = make(map[string]interface{})
			// Override possible values of customData with the specific customData sent
			for k, v := range searchClick.CustomData {
				click.CustomData[k] = v
			}
			if err := click.Execute(v); err != nil {
				return err
			}
		} else {
			return errors.New("Could not find the specific document you are looking for")
		}
	} else {
		Info.Printf("User chose not to click (probability %v%%)", int(searchClick.Probability*100))
	}

	return nil
}
