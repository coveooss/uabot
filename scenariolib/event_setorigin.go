package scenariolib

import "errors"

// ============== SET ORIGIN ======================
// ==================================================

// SetOriginEvent The action of changing the originLevel
type SetOriginEvent struct {
	originLevel1 string
	originLevel2 string
	originLevel3 string
}

func newSetOriginEvent(e *JSONEvent) (*SetOriginEvent, error) {
	var originLevel1, originLevel2, originLevel3 string
	var validCast bool
	if e.Arguments["originLevel1"] != nil {
		if originLevel1, validCast = e.Arguments["originLevel1"].(string); !validCast {
			return nil, errors.New("Parameter originLevel1 must be a string")
		}
	}
	if e.Arguments["originLevel2"] != nil {
		if originLevel2, validCast = e.Arguments["originLevel2"].(string); !validCast {
			return nil, errors.New("Parameter originLevel2 must be a string")
		}
	}
	if e.Arguments["originLevel3"] != nil {
		if originLevel3, validCast = e.Arguments["originLevel3"].(string); !validCast {
			return nil, errors.New("Parameter originLevel3 must be a string")
		}
	}

	return &SetOriginEvent{
		originLevel1: originLevel1,
		originLevel2: originLevel2,
		originLevel3: originLevel3,
	}, nil
}

// Execute Execute the event
func (oe *SetOriginEvent) Execute(v *Visit) error {
	if oe.originLevel1 != "" {
		v.OriginLevel1 = oe.originLevel1
	}
	if oe.originLevel2 != "" {
		v.OriginLevel2 = oe.originLevel2
	}
	if oe.originLevel3 != "" {
		v.OriginLevel3 = oe.originLevel3
	}
	return nil
}
