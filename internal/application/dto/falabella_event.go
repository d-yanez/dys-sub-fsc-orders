package dto

type FalabellaEvent struct {
	Event     string `json:"event"`
	EventType string `json:"eventType"`
	Payload   struct {
		OrderID string `json:"OrderId"`
	} `json:"payload"`
}

func (f FalabellaEvent) EffectiveEventType() string {
	if f.EventType != "" {
		return f.EventType
	}
	return f.Event
}
