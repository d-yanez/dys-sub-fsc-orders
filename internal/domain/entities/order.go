package entities

import "time"

type Order struct {
	ID                   string     `bson:"_id"`
	OrderID              string     `bson:"orderId"`
	OrderNumber          string     `bson:"orderNumber"`
	Status               string     `bson:"status"`
	CreatedAt            *time.Time `bson:"createdAt,omitempty"`
	PromisedShippingTime *time.Time `bson:"promisedShippingTime,omitempty"`
	Marketplace          string     `bson:"marketplace"`
	Audit                AuditOrder `bson:"audit"`
}

type AuditOrder struct {
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
	SourceEvent string    `bson:"sourceEvent"`
}
