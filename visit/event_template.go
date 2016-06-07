package visit

// Template Contains the arguments necessary to send a Template event to the analytics
type Template struct {
	Arg string `json:"arg"`
}

// Parse Parse the different arguments in the JSONEvent to build the event
func (e *Template) Parse(jse *JSONEvent) error {
	return nil
}

// Execute Send the event to the analytics endpoint
func (e *Template) Execute(v *Visit) error {
	return nil
}

// Check for interface implementation
var _ Event = (*Template)(nil)
