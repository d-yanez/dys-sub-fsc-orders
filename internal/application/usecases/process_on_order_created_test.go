package usecases

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
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

type fakeEventLogRepo struct {
	stored      map[string]*entities.EventLog
	markAllowed bool
}

type fakeTelegramNotifier struct {
	sent []ports.TelegramMessage
	err  error
}

func (n *fakeTelegramNotifier) Send(_ context.Context, msg ports.TelegramMessage) error {
	n.sent = append(n.sent, msg)
	return n.err
}

func (r *fakeEventLogRepo) EnsureReceived(_ context.Context, log entities.EventLog) error {
	if r.stored == nil {
		r.stored = map[string]*entities.EventLog{}
	}
	if _, ok := r.stored[log.ID]; !ok {
		copy := log
		r.stored[log.ID] = &copy
	}
	return nil
}

func (r *fakeEventLogRepo) MarkProcessing(_ context.Context, idempotencyKey, messageID, payloadHash, phase string) (bool, error) {
	if r.stored == nil {
		r.stored = map[string]*entities.EventLog{}
	}
	if !r.markAllowed {
		return false, nil
	}
	current, ok := r.stored[idempotencyKey]
	if !ok {
		return false, nil
	}
	if current.Processed {
		return false, nil
	}
	if current.Status != "RECEIVED" && current.Status != "FAILED" {
		return false, nil
	}
	current.Status = "PROCESSING"
	current.Phase = phase
	current.MessageID = messageID
	current.PayloadHash = payloadHash
	current.Attempts += 1
	return true, nil
}

func (r *fakeEventLogRepo) GetByID(_ context.Context, idempotencyKey string) (*entities.EventLog, error) {
	if r.stored == nil {
		return nil, nil
	}
	val, ok := r.stored[idempotencyKey]
	if !ok {
		return nil, nil
	}
	copy := *val
	return &copy, nil
}

func (r *fakeEventLogRepo) MarkCompleted(_ context.Context, idempotencyKey string, status string, warning string) error {
	if r.stored == nil {
		return nil
	}
	if current, ok := r.stored[idempotencyKey]; ok {
		current.Status = status
		current.Processed = true
		current.ErrorSummary = warning
	}
	return nil
}

func (r *fakeEventLogRepo) MarkFailed(_ context.Context, idempotencyKey string, phase string, errSummary string) error {
	if r.stored == nil {
		return nil
	}
	if current, ok := r.stored[idempotencyKey]; ok {
		current.Status = "FAILED"
		current.Processed = false
		current.Phase = phase
		current.ErrorSummary = errSummary
	}
	return nil
}

func TestEvaluateAccepted(t *testing.T) {
	uc := NewProcessOnOrderCreatedUseCase(testLogger(), nil, nil, nil, nil, nil, "")
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
		order:     ports.OrderResponse{OrderID: "1146543495", OrderNumber: "3228563253", Status: "pending"},
		items:     []ports.OrderItemResponse{{OrderItemID: "157246712", SKU: "3516192124", Quantity: 1}},
		thumbnail: &thumb,
	}
	orderRepo := &fakeOrderRepo{}
	itemRepo := &fakeOrderItemRepo{}
	eventRepo := &fakeEventLogRepo{markAllowed: true}
	telegram := &fakeTelegramNotifier{}
	uc := NewProcessOnOrderCreatedUseCase(testLogger(), fscClient, orderRepo, itemRepo, eventRepo, telegram, "https://dy-api-utils-785293986978.us-central1.run.app/stock/view")

	evt := dto.FalabellaEvent{Event: "onOrderCreated"}
	evt.Payload.OrderID = "1146543495"

	result, err := uc.Process(context.Background(), evt, "m-1", "hash")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if result.ItemsCount != 1 {
		t.Fatalf("expected items=1, got %d", result.ItemsCount)
	}
	if orderRepo.last.OrderNumber != "3228563253" {
		t.Fatalf("unexpected order number: %s", orderRepo.last.OrderNumber)
	}
	if len(telegram.sent) != 1 {
		t.Fatalf("expected telegram notification, got %d", len(telegram.sent))
	}
	if !strings.Contains(telegram.sent[0].Text, `ver stock bodega: <a href="https://dy-api-utils-785293986978.us-central1.run.app/stock/view/3516192124">Ver stock bodega</a>`) {
		t.Fatalf("expected stock link in telegram message, got: %s", telegram.sent[0].Text)
	}
}

func TestProcessSuccessTelegramIncludesAllPersistedItems(t *testing.T) {
	thumb := "https://image"
	fscClient := &fakeFSCClient{
		order: ports.OrderResponse{
			OrderID:     "1148513330",
			OrderNumber: "3230404484",
			Status:      "pending",
		},
		items: []ports.OrderItemResponse{
			{OrderItemID: "160707404", SKU: "3737958358", Name: "Item 1", Quantity: 1},
			{OrderItemID: "160707405", SKU: "3737958359", Name: "Item 2", Quantity: 1},
			{OrderItemID: "160707406", SKU: "3737958360", Name: "Item 3", Quantity: 1},
		},
		thumbnail: &thumb,
	}
	orderRepo := &fakeOrderRepo{}
	itemRepo := &fakeOrderItemRepo{}
	eventRepo := &fakeEventLogRepo{markAllowed: true}
	telegram := &fakeTelegramNotifier{}
	uc := NewProcessOnOrderCreatedUseCase(testLogger(), fscClient, orderRepo, itemRepo, eventRepo, telegram, "https://dy-api-utils-785293986978.us-central1.run.app/stock/view")

	evt := dto.FalabellaEvent{Event: "onOrderCreated"}
	evt.Payload.OrderID = "1148513330"

	result, err := uc.Process(context.Background(), evt, "m-3", "hash-3")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if result.ItemsCount != 3 {
		t.Fatalf("expected items=3, got %d", result.ItemsCount)
	}
	if len(telegram.sent) != 3 {
		t.Fatalf("expected 3 telegram notifications (one per item), got %d", len(telegram.sent))
	}

	expectedPerMessage := []struct {
		orderItemID string
		sku         string
		itemName    string
	}{
		{orderItemID: "160707404", sku: "3737958358", itemName: "Item 1"},
		{orderItemID: "160707405", sku: "3737958359", itemName: "Item 2"},
		{orderItemID: "160707406", sku: "3737958360", itemName: "Item 3"},
	}

	for i, sent := range telegram.sent {
		exp := expectedPerMessage[i]
		required := []string{
			"itemsPersistidos: 3",
			"orderItemId: <code>" + exp.orderItemID + "</code>",
			"sku: <code>" + exp.sku + "</code>",
			"item: " + exp.itemName,
			`ver stock bodega: <a href="https://dy-api-utils-785293986978.us-central1.run.app/stock/view/` + exp.sku + `">Ver stock bodega</a>`,
		}
		for _, snippet := range required {
			if !strings.Contains(sent.Text, snippet) {
				t.Fatalf("expected snippet %q in telegram message #%d, got: %s", snippet, i+1, sent.Text)
			}
		}
		if sent.PhotoURL != thumb {
			t.Fatalf("expected photo url %q in telegram message #%d, got %q", thumb, i+1, sent.PhotoURL)
		}
	}
}

func TestProcessSuccessTelegramSendsOneNotificationPerItemForVariableCounts(t *testing.T) {
	cases := []struct {
		name      string
		itemCount int
	}{
		{name: "one_item", itemCount: 1},
		{name: "two_items", itemCount: 2},
		{name: "n_items_five", itemCount: 5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			thumb := "https://image"
			items := make([]ports.OrderItemResponse, 0, tc.itemCount)
			for i := 0; i < tc.itemCount; i++ {
				items = append(items, ports.OrderItemResponse{
					OrderItemID: fmt.Sprintf("160700%03d", i+1),
					SKU:         fmt.Sprintf("373700%03d", i+1),
					Name:        fmt.Sprintf("Item %d", i+1),
					Quantity:    1,
				})
			}

			fscClient := &fakeFSCClient{
				order: ports.OrderResponse{
					OrderID:     "1148513330",
					OrderNumber: "3230404484",
					Status:      "pending",
				},
				items:     items,
				thumbnail: &thumb,
			}
			orderRepo := &fakeOrderRepo{}
			itemRepo := &fakeOrderItemRepo{}
			eventRepo := &fakeEventLogRepo{markAllowed: true}
			telegram := &fakeTelegramNotifier{}
			uc := NewProcessOnOrderCreatedUseCase(testLogger(), fscClient, orderRepo, itemRepo, eventRepo, telegram, "https://dy-api-utils-785293986978.us-central1.run.app/stock/view")

			evt := dto.FalabellaEvent{Event: "onOrderCreated"}
			evt.Payload.OrderID = "1148513330"

			result, err := uc.Process(context.Background(), evt, "m-table", "hash-table")
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if result.ItemsCount != tc.itemCount {
				t.Fatalf("expected items=%d, got %d", tc.itemCount, result.ItemsCount)
			}
			if len(telegram.sent) != tc.itemCount {
				t.Fatalf("expected telegram notifications=%d, got %d", tc.itemCount, len(telegram.sent))
			}

			for i, sent := range telegram.sent {
				expectedOrderItemID := fmt.Sprintf("160700%03d", i+1)
				expectedSKU := fmt.Sprintf("373700%03d", i+1)
				if !strings.Contains(sent.Text, "itemsPersistidos: "+fmt.Sprintf("%d", tc.itemCount)) {
					t.Fatalf("expected itemsPersistidos=%d in message #%d, got: %s", tc.itemCount, i+1, sent.Text)
				}
				if !strings.Contains(sent.Text, "orderItemId: <code>"+expectedOrderItemID+"</code>") {
					t.Fatalf("expected orderItemId %s in message #%d, got: %s", expectedOrderItemID, i+1, sent.Text)
				}
				if !strings.Contains(sent.Text, "sku: <code>"+expectedSKU+"</code>") {
					t.Fatalf("expected sku %s in message #%d, got: %s", expectedSKU, i+1, sent.Text)
				}
			}
		})
	}
}

func TestProcessFailsOnCriticalDependency(t *testing.T) {
	fscClient := &fakeFSCClient{orderErr: errors.New("fsc down")}
	orderRepo := &fakeOrderRepo{}
	itemRepo := &fakeOrderItemRepo{}
	eventRepo := &fakeEventLogRepo{markAllowed: true}
	telegram := &fakeTelegramNotifier{}
	uc := NewProcessOnOrderCreatedUseCase(testLogger(), fscClient, orderRepo, itemRepo, eventRepo, telegram, "")

	evt := dto.FalabellaEvent{Event: "onOrderCreated"}
	evt.Payload.OrderID = "1146543495"

	_, err := uc.Process(context.Background(), evt, "m-1", "hash")
	if err == nil {
		t.Fatal("expected critical error, got nil")
	}
	if len(telegram.sent) != 1 {
		t.Fatalf("expected failed telegram notification, got %d", len(telegram.sent))
	}
}

func TestProcessDuplicateIgnored(t *testing.T) {
	fscClient := &fakeFSCClient{}
	orderRepo := &fakeOrderRepo{}
	itemRepo := &fakeOrderItemRepo{}
	eventRepo := &fakeEventLogRepo{markAllowed: false, stored: map[string]*entities.EventLog{
		"falabella:onOrderCreated:1146543495": {
			ID:        "falabella:onOrderCreated:1146543495",
			Processed: true,
			Status:    "SUCCESS",
		},
	}}
	telegram := &fakeTelegramNotifier{}
	uc := NewProcessOnOrderCreatedUseCase(testLogger(), fscClient, orderRepo, itemRepo, eventRepo, telegram, "")

	evt := dto.FalabellaEvent{Event: "onOrderCreated"}
	evt.Payload.OrderID = "1146543495"

	result, err := uc.Process(context.Background(), evt, "m-2", "hash2")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !result.Duplicate {
		t.Fatal("expected duplicate result")
	}
	if len(telegram.sent) != 0 {
		t.Fatalf("expected no telegram for duplicate, got %d", len(telegram.sent))
	}
}

func TestBuildTelegramMessageWithoutSKUDoesNotIncludeStockLink(t *testing.T) {
	msg := buildTelegramMessage(ProcessResult{
		Status:    "SUCCESS",
		EventType: "onOrderCreated",
		OrderID:   "1146543495",
		FirstSKU:  "",
	}, "", "https://dy-api-utils-785293986978.us-central1.run.app/stock/view")

	if strings.Contains(msg.Text, "Ver stock bodega") {
		t.Fatalf("did not expect stock link when sku is empty, got: %s", msg.Text)
	}
}

func TestBuildTelegramMessageEscapesSKUInStockLink(t *testing.T) {
	msg := buildTelegramMessage(ProcessResult{
		Status:    "SUCCESS",
		EventType: "onOrderCreated",
		OrderID:   "1146543495",
		FirstSKU:  "ABC/123",
	}, "", "https://dy-api-utils-785293986978.us-central1.run.app/stock/view/")

	if !strings.Contains(msg.Text, `href="https://dy-api-utils-785293986978.us-central1.run.app/stock/view/ABC%2F123"`) {
		t.Fatalf("expected escaped sku in stock link, got: %s", msg.Text)
	}
}
