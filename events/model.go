package events

/*
Event - represents the event
*/
type Event struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}
