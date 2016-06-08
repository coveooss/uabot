package scenariolib

import "errors"

// ============== VIEW EVENT ======================
// ==================================================

// ViewEvent a struct representing a search, is defined by a query to execute
type ViewEvent struct {
	pageURI         string
	pageReferrer    string
	pageTitle       string
	contentIdKey    string
	contentIdValue  string
	contentType     string
}

func newViewEvent(e *JSONEvent) (*ViewEvent, error) {
	pageURI, ok1 := e.Arguments["pageuri"].(string)
	pageReferrer, ok2 := e.Arguments["pagereferrer"].(string)
	pageTitle, ok3 := e.Arguments["pagetitle"].(string)
	contentIdKey, ok4 := e.Arguments["contentidkey"].(string)
	contentIdValue, ok5 := e.Arguments["contentidvalue"].(string)
	contentType, ok6 := e.Arguments["contenttype"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return nil, errors.New("Invalid parse of arguments on View Event")
	}

	if pageURI == "" {
		// Find a random page to view
	}

	return &ViewEvent{
		pageURI:      	pageURI,
		pageReferrer: 	pageReferrer,
		pageTitle:    	pageTitle,
		contentIdKey: 	contentIdKey,
		contentIdValue: contentIdValue,
		contentType:  	contentType,
	}, nil
}

// Execute Execute the search event, runs the query and sends a search event to
// the analytics.
func (ve *ViewEvent) Execute(v *Visit) error {

	err := v.sendViewEvent(ve.pageTitle, ve.pageReferrer, ve.pageURI, ve.contentIdKey, ve.contentIdValue, ve.contentType)

	return err
}
