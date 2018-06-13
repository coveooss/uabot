// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

// ============== TAB CHANGE EVENT ======================
// ======================================================

// TabChangeEvent represents a tab change event
type TabChangeEvent struct {
	Name               string `json:"name"`
	ConstantExpression string `json:"cq,omitempty"`
}

// IsValid Additional validation after the json unmarshal.
func (tab *TabChangeEvent) IsValid() (bool, string) {
	return true, ""
}

// Execute Sends the tabchange event to the analytics and modify the CQ for the
// following queries in the visit
func (tab *TabChangeEvent) Execute(v *Visit) error {
	Info.Printf("Changing tab to %s with CQ : %s", tab.Name, tab.ConstantExpression)

	v.LastQuery.CQ = v.LastQuery.CQ + " " + tab.ConstantExpression
	v.OriginLevel2 = tab.Name
	v.LastQuery.Tab = tab.Name

	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}
	v.LastResponse = resp

	Info.Printf("Sending TabChange Event : %s", tab.Name)
	err = v.sendInterfaceChangeEvent("interfaceChange", "", map[string]interface{}{"interfaceChangeTo": v.OriginLevel2})
	if err != nil {
		return err
	}
	return nil
}
