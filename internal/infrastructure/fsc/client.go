package fsc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/ports"
)

type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func NewClient(baseURL, apiKey string, timeout time.Duration) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GetOrder(ctx context.Context, orderID string) (ports.OrderResponse, error) {
	body, err := c.getJSON(ctx, "/order/"+url.PathEscape(orderID))
	if err != nil {
		return ports.OrderResponse{}, err
	}

	financial := firstMap(body, "financial", "Financial")
	totals := firstMap(body, "totals", "Totals")
	grandTotal := pickNumber(
		firstNumber(financial, "grandTotal", "GrandTotal"),
		firstNumber(totals, "grandTotal", "GrandTotal"),
		firstNumber(body, "grandTotal", "GrandTotal", "Price"),
	)
	productTotal := pickNumber(
		firstNumber(financial, "productTotal", "ProductTotal"),
		firstNumber(totals, "productTotal", "ProductTotal"),
		firstNumber(body, "productTotal", "ProductTotal"),
	)
	taxAmount := pickNumber(
		firstNumber(financial, "taxAmount", "TaxAmount"),
		firstNumber(totals, "taxAmount", "TaxAmount"),
		firstNumber(body, "taxAmount", "TaxAmount"),
	)
	shippingFeeTotal := pickNumber(
		firstNumber(financial, "shippingFeeTotal", "ShippingFeeTotal"),
		firstNumber(totals, "shippingFeeTotal", "ShippingFeeTotal"),
		firstNumber(body, "shippingFeeTotal", "ShippingFeeTotal"),
	)

	return ports.OrderResponse{
		OrderID:              firstString(body, "orderId", "OrderId", "id"),
		OrderNumber:          firstString(body, "orderNumber", "OrderNumber"),
		Status:               firstString(body, "status", "Status"),
		CreatedAt:            firstString(body, "createdAt", "CreatedAt"),
		PromisedShippingTime: firstString(body, "promisedShippingTime", "PromisedShippingTime"),
		InvoiceRequired:      firstBool(body, "invoiceRequired", "InvoiceRequired"),
		GrandTotal:           grandTotal,
		ProductTotal:         productTotal,
		TaxAmount:            taxAmount,
		ShippingFeeTotal:     shippingFeeTotal,
		AddressBilling:       firstMap(body, "addressBilling", "AddressBilling"),
		AddressShipping:      firstMap(body, "addressShipping", "AddressShipping"),
	}, nil
}

func (c *Client) GetOrderItems(ctx context.Context, orderID string) ([]ports.OrderItemResponse, error) {
	body, err := c.getJSON(ctx, "/orderItems/"+url.PathEscape(orderID))
	if err != nil {
		return nil, err
	}

	itemsRaw, ok := body["items"].([]any)
	if !ok {
		return []ports.OrderItemResponse{}, nil
	}

	items := make([]ports.OrderItemResponse, 0, len(itemsRaw))
	for _, it := range itemsRaw {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		item := ports.OrderItemResponse{
			OrderItemID:      firstString(m, "orderItemId", "OrderItemId", "id"),
			Name:             firstString(m, "name", "Name"),
			SKU:              firstString(m, "sku", "Sku", "SKU"),
			ShopSKU:          firstString(m, "shopSku", "ShopSku"),
			Status:           firstString(m, "status", "Status"),
			Quantity:         firstInt(m, "quantity", "Quantity"),
			Price:            firstInt64(m, "price", "Price", "itemPrice", "ItemPrice"),
			PaidPrice:        firstInt64(m, "paidPrice", "PaidPrice"),
			ShippingAmount:   firstInt64(m, "shippingAmount", "ShippingAmount"),
			TrackingCode:     firstString(m, "trackingCode", "TrackingCode"),
			PackageID:        firstString(m, "packageId", "PackageId", "package_id"),
			ShipmentProvider: firstString(m, "shipmentProvider", "ShipmentProvider"),
		}
		items = append(items, item)
	}

	return items, nil
}

func (c *Client) GetSKUThumbnail(ctx context.Context, sku string) (*string, error) {
	path := "/sku/" + url.PathEscape(sku) + "?fields=sku,name,price,stock,quantity,images,status"
	body, err := c.getJSON(ctx, path)
	if err != nil {
		return nil, err
	}

	imagesRaw, ok := body["images"].([]any)
	if !ok || len(imagesRaw) == 0 {
		return nil, nil
	}

	switch v := imagesRaw[0].(type) {
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return nil, nil
		}
		return &trimmed, nil
	case map[string]any:
		candidate := firstString(v, "url", "secure_url", "src")
		if candidate == "" {
			return nil, nil
		}
		return &candidate, nil
	default:
		return nil, nil
	}
}

func (c *Client) getJSON(ctx context.Context, path string) (map[string]any, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("fsc http error status=%d path=%s body=%s", resp.StatusCode, path, string(bodyBytes))
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func pickNumber(values ...*float64) *float64 {
	for _, v := range values {
		if v != nil {
			return v
		}
	}
	return nil
}
