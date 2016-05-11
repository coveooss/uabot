package scenariolib

import "errors"

// ============== VIEW EVENT ======================
// ==================================================

// ViewEvent a struct representing a search, is defined by a query to execute
type ViewEvent struct {
	pageURI      string
	pageReferrer string
	pageTitle    string
}

func newViewEvent(e *JSONEvent) (*ViewEvent, error) {
	pageURI, ok1 := e.Arguments["pageuri"].(string)
	pageReferrer, ok2 := e.Arguments["pagereferrer"].(string)
	pageTitle, ok3 := e.Arguments["pagetitle"].(string)
	if !ok1 || !ok2 || !ok3 {
		return nil, errors.New("ERR >>> Invalid parse of arguments on View Event")
	}

	if pageURI == "" {
		// Find a random page to view
	}

	return &ViewEvent{
		pageURI:      pageURI,
		pageReferrer: pageReferrer,
		pageTitle:    pageTitle,
	}, nil
}

// Execute Execute the search event, runs the query and sends a search event to
// the analytics.
func (ve *ViewEvent) Execute(v *Visit) error {

	err := v.sendViewEvent(ve.pageTitle, ve.pageReferrer, ve.pageURI)

	return err
}
