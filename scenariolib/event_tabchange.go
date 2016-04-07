// Package scenariolib handles everything need to execute a scenario and send all
// information to the usage analytics endpoint
package scenariolib

import (
	"errors"

	"github.com/k0kubun/pp"
)

// ============== TAB CHANGE EVENT ======================
// ======================================================

// TabChangeEvent represents a tab change event
type TabChangeEvent struct {
	name string
	cq   string
}

func newTabChangeEvent(e *JSONEvent) (*TabChangeEvent, error) {
	name, ok1 := e.Arguments["tabName"].(string)
	cq, ok2 := e.Arguments["tabCQ"].(string)
	if !ok1 || !ok2 {
		return nil, errors.New("ERR >>> Invalid parse of arguments on TabChange Event")
	}

	return &TabChangeEvent{
		name: name,
		cq:   cq,
	}, nil
}

// Execute Sends the tabchange event to the analytics and modify the CQ for the
// following queries in the visit
func (tc *TabChangeEvent) Execute(v *Visit) error {
	pp.Printf("\nLOG >>> Changing tab to %v with CQ : %v", tc.name, tc.cq)

	v.LastQuery.CQ = tc.cq
	v.OriginLevel2 = tc.name
	v.LastQuery.Tab = tc.name

	resp, err := v.SearchClient.Query(*v.LastQuery)
	if err != nil {
		return err
	}
	v.LastResponse = resp

	pp.Printf("\nLOG >>> Sending Tab Change Event : %v", tc.name)
	err = v.sendInterfaceChangeEvent()
	if err != nil {
		return err
	}
	return nil
}

func (tc *TabChangeEvent) parseTabChangeEvent(e *JSONEvent) error {
	name, ok1 := e.Arguments["tabName"].(string)
	cq, ok2 := e.Arguments["tabCQ"].(string)
	if !ok1 || !ok2 {
		return errors.New("ERR >>> Invalid parse of arguments on TabChange Event")
	}

	tc.name = name
	tc.cq = cq

	return nil
}
