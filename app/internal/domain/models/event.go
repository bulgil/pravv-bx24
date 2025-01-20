package models

type Event struct {
	EventType EventType `json:"event"`
	Data      Data      `json:"data"`
}

func (e Event) IsNil() bool {
	return e.EventType == ""
}
