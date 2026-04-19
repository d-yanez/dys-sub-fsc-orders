package entities

import "time"

type Order struct {
	ID                   string     `bson:"_id"`
	OrderID              string     `bson:"orderId"`
	OrderNumber          string     `bson:"orderNumber"`
	Status               string     `bson:"status"`
	CreatedAt            *time.Time `bson:"createdAt,omitempty"`
	PromisedShippingTime *time.Time `bson:"promisedShippingTime,omitempty"`
	Financial            Financial  `bson:"financial"`
	Addresses            Addresses  `bson:"addresses"`
	Marketplace          string     `bson:"marketplace"`
	Audit                AuditOrder `bson:"audit"`
}

type Financial struct {
	GrandTotal       *float64 `bson:"grandTotal"`
	ProductTotal     *float64 `bson:"productTotal"`
	TaxAmount        *float64 `bson:"taxAmount"`
	ShippingFeeTotal *float64 `bson:"shippingFeeTotal"`
	InvoiceRequired  *bool    `bson:"invoiceRequired"`
	DocumentType     *string  `bson:"documentType"`
}

type Addresses struct {
	Billing  map[string]any `bson:"billing"`
	Shipping map[string]any `bson:"shipping"`
}

type AuditOrder struct {
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
	SourceEvent string    `bson:"sourceEvent"`
}
