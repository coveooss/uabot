package scenariolib

import "errors"

// ============== SET ORIGIN ======================
// ==================================================

// SetOriginEvent The action of changing the originLevel
type SetOriginEvent struct {
	originLevel1          string
	originLevel2          string
	originLevel3          string
	randomizeOriginLevel1 bool
	randomizeOriginLevel2 bool
}

func newSetOriginEvent(e *JSONEvent) (*SetOriginEvent, error) {
	var originLevel1, originLevel2, originLevel3 string
	var validCast bool
	var randomizeOriginLevel1, randomizeOriginLevel2 bool

	if e.Arguments["randomizeOriginLevel1"] != nil {
		if randomizeOriginLevel1, validCast = e.Arguments["randomizeOriginLevel1"].(bool); !validCast {
			return nil, errors.New("Parameter randomzieOriginLevel1 must be a boolean")
		}
	}
	if e.Arguments["randomizeOriginLevel2"] != nil {
		if randomizeOriginLevel2, validCast = e.Arguments["randomizeOriginLevel2"].(bool); !validCast {
			return nil, errors.New("Parameter randomzieOriginLevel2 must be a boolean")
		}
	}

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
		originLevel1:          originLevel1,
		originLevel2:          originLevel2,
		originLevel3:          originLevel3,
		randomizeOriginLevel1: randomizeOriginLevel1,
		randomizeOriginLevel2: randomizeOriginLevel2,
	}, nil
}

// Execute Execute the event
func (oe *SetOriginEvent) Execute(v *Visit) error {
	var err error
	if oe.randomizeOriginLevel1 {
		v.OriginLevel1, err = v.Config.RandomOriginLevel1()
		if err != nil {
			return err
		}
	} else {
		if oe.originLevel1 != "" {
			v.OriginLevel1 = oe.originLevel1
		}
	}
	if oe.randomizeOriginLevel2 {
		v.OriginLevel2, err = v.Config.RandomOriginLevel2(v.OriginLevel1)
		if err != nil {
			return err
		}
	} else {
		if oe.originLevel2 != "" {
			v.OriginLevel2 = oe.originLevel2
		}
	}
	if oe.originLevel3 != "" {
		v.OriginLevel3 = oe.originLevel3
	}

	return nil
}
