package ports

import "context"

type FSCClient interface {
	GetOrder(ctx context.Context, orderID string) (OrderResponse, error)
	GetOrderItems(ctx context.Context, orderID string) ([]OrderItemResponse, error)
	GetSKUThumbnail(ctx context.Context, sku string) (*string, error)
}

type OrderResponse struct {
	OrderID              string
	OrderNumber          string
	Status               string
	CreatedAt            string
	PromisedShippingTime string
}

type OrderItemResponse struct {
	OrderItemID      string
	Name             string
	SKU              string
	ShopSKU          string
	Status           string
	Quantity         int
	TrackingCode     string
	PackageID        string
	ShipmentProvider string
}
