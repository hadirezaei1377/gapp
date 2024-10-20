package entity

type Notification struct {
	EventType string `json:"event_type"`
	Payload   string `json:"payload"`
}
