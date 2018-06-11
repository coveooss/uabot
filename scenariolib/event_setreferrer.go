package scenariolib

// ============== SET REFERRER ======================
// ==================================================

// SetReferrerEvent The action of changing the referrer
type SetReferrerEvent struct {
	Referrer string `json:"referrer"`
}

// IsValid Additional validation after the json unmarshal.
func (referrer *SetReferrerEvent) IsValid() (bool, string) {
	return true, ""
}

//Execute Execute the event
func (referrer *SetReferrerEvent) Execute(v *Visit) error {
	v.Referrer = referrer.Referrer
	return nil
}
