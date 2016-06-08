package refactor

import "errors"

// Event Generic interface for events to implement
type Event interface {
	Parse(*JSONEvent) error
	Execute(*Visit) error
}

// CreateEvent Create an event from a JSONEvent
func CreateEvent(e *JSONEvent) (event Event, err error) {
	switch e.Type {
	case "Click":
		event = new(ClickEvent)
		err = event.Parse(e)
	default:
		return nil, errors.New("Unsupported type of events")
	}

	return event, err
}
