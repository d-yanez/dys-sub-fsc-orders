package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/application/usecases"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/entities"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/ports"
)

func handlerTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type fakeFSCClient struct{}

func (f *fakeFSCClient) GetOrder(_ context.Context, orderID string) (ports.OrderResponse, error) {
	return ports.OrderResponse{OrderID: orderID, OrderNumber: "3228563253", Status: "pending"}, nil
}

func (f *fakeFSCClient) GetOrderItems(_ context.Context, _ string) ([]ports.OrderItemResponse, error) {
	return []ports.OrderItemResponse{{OrderItemID: "157246712", SKU: "3516192124", Quantity: 1}}, nil
}

func (f *fakeFSCClient) GetSKUThumbnail(_ context.Context, _ string) (*string, error) {
	url := "https://example.com/image.jpg"
	return &url, nil
}

type fakeOrderRepo struct{}

func (r *fakeOrderRepo) Upsert(_ context.Context, _ entities.Order) error {
	return nil
}

type fakeOrderItemRepo struct{}

func (r *fakeOrderItemRepo) UpsertMany(_ context.Context, _ []entities.OrderItem) error {
	return nil
}

func TestHandleAcceptedOnOrderCreated(t *testing.T) {
	uc := usecases.NewProcessOnOrderCreatedUseCase(handlerTestLogger(), &fakeFSCClient{}, &fakeOrderRepo{}, &fakeOrderItemRepo{})
	h := NewPubSubPushHandler(handlerTestLogger(), uc)

	eventJSON := `{"event":"onOrderCreated","payload":{"OrderId":"1146543495"}}`
	data := base64.StdEncoding.EncodeToString([]byte(eventJSON))
	body := `{"message":{"messageId":"m-1","data":"` + data + `"}}`

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"status":"SUCCESS"`)) {
		t.Fatalf("expected SUCCESS response, got body=%s", rec.Body.String())
	}
}

func TestHandleIgnoredMissingData(t *testing.T) {
	uc := usecases.NewProcessOnOrderCreatedUseCase(handlerTestLogger(), &fakeFSCClient{}, &fakeOrderRepo{}, &fakeOrderItemRepo{})
	h := NewPubSubPushHandler(handlerTestLogger(), uc)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"message":{"messageId":"m-1"}}`))
	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"status":"ignored"`)) {
		t.Fatalf("expected ignored response, got body=%s", rec.Body.String())
	}
}
