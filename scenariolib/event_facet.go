package scenariolib

import (
	"errors"
	"fmt"
)

// ============== FACET CHANGE EVENT ======================
// ======================================================

// FacetEvent represents a tab change event
type FacetEvent struct {
	Title string
	Value string
	Field string
}

func newFacetEvent(e *JSONEvent) (*FacetEvent, error) {
	Title, ok1 := e.Arguments["facetTitle"].(string)
	Value, ok2 := e.Arguments["facetValue"].(string)
	Field, ok3 := e.Arguments["facetField"].(string)
	if !ok1 || !ok2 || !ok3 {
		return nil, errors.New("Invalid parse of arguments on Facet Event")
	}

	return &FacetEvent{
		Title: Title,
		Value: Value,
		Field: Field,
	}, nil
}

// Execute Sends the tabchange event to the analytics and modify the CQ for the
// following queries in the visit
func (e *FacetEvent) Execute(v *Visit) error {
	Info.Printf("Clicking on facet title=%s value=%s", e.Title, e.Value)

	v.LastQuery.AQ = fmt.Sprintf("%s==\"%s\"", e.Field, e.Value)

	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}
	v.LastResponse = resp

	Info.Printf("Sending FacetChange Event title=%s value=%s", e.Title, e.Value)
	customData := make(map[string]interface{})
	customData["facetValue"] = e.Value
	customData["facetTitle"] = e.Title
	customData["facetField"] = e.Field
	err = v.sendInterfaceChangeEvent("facetSelect", "facet", customData)
	if err != nil {
		return err
	}
	return nil
}
