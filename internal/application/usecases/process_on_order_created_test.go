package usecases

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/application/dto"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/entities"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/ports"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type fakeFSCClient struct {
	order        ports.OrderResponse
	items        []ports.OrderItemResponse
	thumbnail    *string
	orderErr     error
	itemsErr     error
	thumbnailErr error
}

func (f *fakeFSCClient) GetOrder(_ context.Context, _ string) (ports.OrderResponse, error) {
	if f.orderErr != nil {
		return ports.OrderResponse{}, f.orderErr
	}
	return f.order, nil
}

func (f *fakeFSCClient) GetOrderItems(_ context.Context, _ string) ([]ports.OrderItemResponse, error) {
	if f.itemsErr != nil {
		return nil, f.itemsErr
	}
	return f.items, nil
}

func (f *fakeFSCClient) GetSKUThumbnail(_ context.Context, _ string) (*string, error) {
	if f.thumbnailErr != nil {
		return nil, f.thumbnailErr
	}
	return f.thumbnail, nil
}

type fakeOrderRepo struct {
	last entities.Order
	err  error
}

func (r *fakeOrderRepo) Upsert(_ context.Context, order entities.Order) error {
	if r.err != nil {
		return r.err
	}
	r.last = order
	return nil
}

type fakeOrderItemRepo struct {
	last []entities.OrderItem
	err  error
}

func (r *fakeOrderItemRepo) UpsertMany(_ context.Context, items []entities.OrderItem) error {
	if r.err != nil {
		return r.err
	}
	r.last = items
	return nil
}

func TestEvaluateAccepted(t *testing.T) {
	uc := NewProcessOnOrderCreatedUseCase(testLogger(), nil, nil, nil)
	evt := dto.FalabellaEvent{Event: "onOrderCreated"}
	evt.Payload.OrderID = "1146543495"

	got := uc.Evaluate(evt)
	if !got.ShouldProcess {
		t.Fatalf("expected ShouldProcess=true, got false reason=%s", got.Reason)
	}
}

func TestProcessSuccess(t *testing.T) {
	thumb := "https://image"
	fscClient := &fakeFSCClient{
		order: ports.OrderResponse{OrderID: "1146543495", OrderNumber: "3228563253", Status: "pending"},
		items: []ports.OrderItemResponse{
			{OrderItemID: "157246712", SKU: "3516192124", Quantity: 1},
		},
		thumbnail: &thumb,
	}
	orderRepo := &fakeOrderRepo{}
	itemRepo := &fakeOrderItemRepo{}
	uc := NewProcessOnOrderCreatedUseCase(testLogger(), fscClient, orderRepo, itemRepo)

	evt := dto.FalabellaEvent{Event: "onOrderCreated"}
	evt.Payload.OrderID = "1146543495"

	result, err := uc.Process(context.Background(), evt, "m-1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if result.ItemsCount != 1 {
		t.Fatalf("expected items=1, got %d", result.ItemsCount)
	}
	if orderRepo.last.OrderNumber != "3228563253" {
		t.Fatalf("unexpected order number: %s", orderRepo.last.OrderNumber)
	}
}

func TestProcessFailsOnCriticalDependency(t *testing.T) {
	fscClient := &fakeFSCClient{orderErr: errors.New("fsc down")}
	orderRepo := &fakeOrderRepo{}
	itemRepo := &fakeOrderItemRepo{}
	uc := NewProcessOnOrderCreatedUseCase(testLogger(), fscClient, orderRepo, itemRepo)

	evt := dto.FalabellaEvent{Event: "onOrderCreated"}
	evt.Payload.OrderID = "1146543495"

	_, err := uc.Process(context.Background(), evt, "m-1")
	if err == nil {
		t.Fatal("expected critical error, got nil")
	}
}
