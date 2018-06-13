package scenariolib

import (
	"fmt"
)

// ============== FACET CHANGE EVENT ======================
// ======================================================

// FacetEvent represents a tab change event
type FacetEvent struct {
	FacetTitle string                 `json:"facetTitle"`
	FacetValue string                 `json:"facetValue"`
	FacetField string                 `json:"facetField"`
	CustomData map[string]interface{} `json:"customData,omitempty"`
}

// IsValid Additional validation after the json unmarshal.
func (facet *FacetEvent) IsValid() (bool, string) {
	return true, ""
}

// Execute Sends the tabchange event to the analytics and modify the CQ for the
// following queries in the visit
func (facet *FacetEvent) Execute(v *Visit) error {
	Info.Printf("Clicking on facet title=%s value=%s", facet.FacetTitle, facet.FacetValue)

	v.LastQuery.AQ = fmt.Sprintf("%s==\"%s\"", facet.FacetField, facet.FacetValue)

	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}
	v.LastResponse = resp

	Info.Printf("Sending FacetChange Event title=%s value=%s", facet.FacetTitle, facet.FacetValue)
	if facet.CustomData == nil {
		facet.CustomData = make(map[string]interface{})
	}
	facet.CustomData["facetValue"] = facet.FacetValue
	facet.CustomData["facetTitle"] = facet.FacetTitle
	facet.CustomData["facetId"] = facet.FacetField
	err = v.sendInterfaceChangeEvent("facetSelect", "facet", facet.CustomData)
	if err != nil {
		return err
	}
	return nil
}
