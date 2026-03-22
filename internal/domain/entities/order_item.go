package entities

import "time"

type OrderItem struct {
	ID               string         `bson:"_id"`
	OrderItemID      string         `bson:"orderItemId"`
	OrderID          string         `bson:"orderId"`
	OrderNumber      string         `bson:"orderNumber"`
	Name             string         `bson:"name,omitempty"`
	SKU              string         `bson:"sku,omitempty"`
	ShopSKU          string         `bson:"shopSku,omitempty"`
	Status           string         `bson:"status,omitempty"`
	Quantity         int            `bson:"quantity"`
	TrackingCode     string         `bson:"trackingCode,omitempty"`
	PackageID        string         `bson:"packageId,omitempty"`
	ShipmentProvider string         `bson:"shipmentProvider,omitempty"`
	Thumbnail        *string        `bson:"thumbnail,omitempty"`
	Audit            AuditOrderItem `bson:"audit"`
}

type AuditOrderItem struct {
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}
