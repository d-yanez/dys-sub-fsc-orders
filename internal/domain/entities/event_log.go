package entities

import "time"

type EventLog struct {
	ID             string     `bson:"_id"`
	IdempotencyKey string     `bson:"idempotencyKey"`
	MessageID      string     `bson:"messageId"`
	EventType      string     `bson:"eventType"`
	OrderID        string     `bson:"orderId"`
	PayloadHash    string     `bson:"payloadHash"`
	Status         string     `bson:"status"`
	Processed      bool       `bson:"processed"`
	Attempts       int        `bson:"attempts"`
	Phase          string     `bson:"phase"`
	ProcessedAt    *time.Time `bson:"processedAt,omitempty"`
	ErrorSummary   string     `bson:"errorSummary,omitempty"`
	CreatedAt      time.Time  `bson:"createdAt"`
	UpdatedAt      time.Time  `bson:"updatedAt"`
}
