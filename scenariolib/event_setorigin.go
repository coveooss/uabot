package scenariolib

// ============== SET ORIGIN ======================
// ==================================================

// SetOriginEvent The action of changing the originLevel1-2-3 of the current visit.
type SetOriginEvent struct {
	OriginLevel1 string `json:"originLevel1"`
	OriginLevel2 string `json:"originLevel2,omitempty"`
	OriginLevel3 string `json:"originLevel3,omitempty"`
}

// IsValid Additional validation after the json unmarshal.
func (origin *SetOriginEvent) IsValid() (bool, string) {
	return true, ""
}

// Execute the set origin event. Replaces the originLevel1-2-3 in the current visit.
func (origin *SetOriginEvent) Execute(v *Visit) error {
	Info.Printf("Executing SetOrigin {originLevel1: %s, originLevel2: %s, OriginLevel3: %s}", origin.OriginLevel1, origin.OriginLevel2, origin.OriginLevel3)
	if origin.OriginLevel1 != "" {
		v.OriginLevel1 = origin.OriginLevel1
	}
	if origin.OriginLevel2 != "" {
		v.OriginLevel2 = origin.OriginLevel2
	}
	if origin.OriginLevel3 != "" {
		v.OriginLevel3 = origin.OriginLevel3
	}
	return nil
}
