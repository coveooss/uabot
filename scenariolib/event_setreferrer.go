package scenariolib

import "errors"

// ============== SET REFERRER ======================
// ==================================================

// SetReferrerEvent The action of changing the referrer
type SetReferrerEvent struct {
	referrer string
}

func newSetReferrerEvent(e *JSONEvent) (*SetReferrerEvent, error) {
	var referrer string
	var validCast bool
	if e.Arguments["referrer"] != nil {
		if referrer, validCast = e.Arguments["referrer"].(string); !validCast {
			return nil, errors.New("Parameter referrer must be a string")
		}
	}

	return &SetReferrerEvent{
		referrer: referrer,
	}, nil
}

//Execute Execute the event
func (oe *SetReferrerEvent) Execute(v *Visit) error {
	if oe.referrer != "" {
		v.Referrer = oe.referrer
	}
	return nil
}
