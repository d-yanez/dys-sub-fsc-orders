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
	Customer             Customer   `bson:"customer,omitempty"`
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

type Customer struct {
	FirstName                  string         `bson:"firstName,omitempty"`
	LastName                   string         `bson:"lastName,omitempty"`
	Email                      string         `bson:"email,omitempty"`
	NationalRegistrationNumber string         `bson:"nationalRegistrationNumber,omitempty"`
	Company                    string         `bson:"company,omitempty"`
	Activity                   string         `bson:"activity,omitempty"`
	Address                    string         `bson:"address,omitempty"`
	Municipality               string         `bson:"municipality,omitempty"`
	ExtraBillingAttributes     map[string]any `bson:"extraBillingAttributes,omitempty"`
}

type AuditOrder struct {
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
	SourceEvent string    `bson:"sourceEvent"`
}
