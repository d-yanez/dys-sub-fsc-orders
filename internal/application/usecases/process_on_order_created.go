package usecases

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/application/dto"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/entities"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/enums"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/ports"
)

type ProcessDecision struct {
	ShouldProcess bool
	EventType     string
	OrderID       string
	Reason        string
}

type ProcessOnOrderCreatedUseCase struct {
	log           *slog.Logger
	fscClient     ports.FSCClient
	orderRepo     ports.OrderRepository
	orderItemRepo ports.OrderItemRepository
}

type ProcessResult struct {
	Status      enums.ResultStatus
	EventType   string
	OrderID     string
	OrderNumber string
	ItemsCount  int
	Warning     string
}

func NewProcessOnOrderCreatedUseCase(
	log *slog.Logger,
	fscClient ports.FSCClient,
	orderRepo ports.OrderRepository,
	orderItemRepo ports.OrderItemRepository,
) *ProcessOnOrderCreatedUseCase {
	return &ProcessOnOrderCreatedUseCase{
		log:           log,
		fscClient:     fscClient,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

func (u *ProcessOnOrderCreatedUseCase) Evaluate(event dto.FalabellaEvent) ProcessDecision {
	eventType := strings.TrimSpace(event.EffectiveEventType())
	orderID := strings.TrimSpace(event.Payload.OrderID)

	if eventType == "" {
		return ProcessDecision{ShouldProcess: false, EventType: eventType, OrderID: orderID, Reason: "missing_event_type"}
	}
	if eventType != "onOrderCreated" {
		return ProcessDecision{ShouldProcess: false, EventType: eventType, OrderID: orderID, Reason: "unsupported_event_type"}
	}
	if orderID == "" {
		return ProcessDecision{ShouldProcess: false, EventType: eventType, OrderID: orderID, Reason: "missing_order_id"}
	}

	return ProcessDecision{ShouldProcess: true, EventType: eventType, OrderID: orderID, Reason: "accepted_for_processing"}
}

func (u *ProcessOnOrderCreatedUseCase) Process(ctx context.Context, event dto.FalabellaEvent, messageID string) (ProcessResult, error) {
	decision := u.Evaluate(event)
	if !decision.ShouldProcess {
		return ProcessResult{
			Status:    enums.ResultSuccess,
			EventType: decision.EventType,
			OrderID:   decision.OrderID,
			Warning:   decision.Reason,
		}, nil
	}

	if u.fscClient == nil || u.orderRepo == nil || u.orderItemRepo == nil {
		return ProcessResult{}, errors.New("usecase dependencies are not configured")
	}

	orderResp, err := u.fscClient.GetOrder(ctx, decision.OrderID)
	if err != nil {
		return ProcessResult{}, err
	}
	itemsResp, err := u.fscClient.GetOrderItems(ctx, decision.OrderID)
	if err != nil {
		return ProcessResult{}, err
	}

	now := time.Now().UTC()
	orderCreatedAt := parseTime(orderResp.CreatedAt)
	promisedShippingTime := parseTime(orderResp.PromisedShippingTime)

	order := entities.Order{
		ID:                   decision.OrderID,
		OrderID:              decision.OrderID,
		OrderNumber:          strings.TrimSpace(orderResp.OrderNumber),
		Status:               strings.TrimSpace(orderResp.Status),
		CreatedAt:            orderCreatedAt,
		PromisedShippingTime: promisedShippingTime,
		Marketplace:          "falabella",
		Audit: entities.AuditOrder{
			CreatedAt:   now,
			UpdatedAt:   now,
			SourceEvent: decision.EventType,
		},
	}

	if err := u.orderRepo.Upsert(ctx, order); err != nil {
		return ProcessResult{}, err
	}

	partialWarning := ""
	orderItems := make([]entities.OrderItem, 0, len(itemsResp))
	for _, it := range itemsResp {
		itemID := strings.TrimSpace(it.OrderItemID)
		if itemID == "" {
			continue
		}

		var thumbnail *string
		sku := strings.TrimSpace(it.SKU)
		if sku != "" {
			result, thumbErr := u.fscClient.GetSKUThumbnail(ctx, sku)
			if thumbErr != nil && partialWarning == "" {
				partialWarning = "thumbnail_lookup_failed"
				u.log.Warn("thumbnail_lookup_failed",
					"orderId", decision.OrderID,
					"sku", sku,
					"error", thumbErr.Error(),
					"messageId", messageID,
				)
			}
			thumbnail = result
		}

		orderItems = append(orderItems, entities.OrderItem{
			ID:               itemID,
			OrderItemID:      itemID,
			OrderID:          decision.OrderID,
			OrderNumber:      order.OrderNumber,
			Name:             strings.TrimSpace(it.Name),
			SKU:              sku,
			ShopSKU:          strings.TrimSpace(it.ShopSKU),
			Status:           strings.TrimSpace(it.Status),
			Quantity:         it.Quantity,
			TrackingCode:     strings.TrimSpace(it.TrackingCode),
			PackageID:        strings.TrimSpace(it.PackageID),
			ShipmentProvider: strings.TrimSpace(it.ShipmentProvider),
			Thumbnail:        thumbnail,
			Audit: entities.AuditOrderItem{
				CreatedAt: now,
				UpdatedAt: now,
			},
		})
	}

	if err := u.orderItemRepo.UpsertMany(ctx, orderItems); err != nil {
		return ProcessResult{}, err
	}

	status := enums.ResultSuccess
	if partialWarning != "" {
		status = enums.ResultPartialSuccess
	}

	return ProcessResult{
		Status:      status,
		EventType:   decision.EventType,
		OrderID:     decision.OrderID,
		OrderNumber: order.OrderNumber,
		ItemsCount:  len(orderItems),
		Warning:     partialWarning,
	}, nil
}

func parseTime(raw string) *time.Time {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil
	}

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		parsed, err := time.Parse(format, value)
		if err == nil {
			return &parsed
		}
	}
	return nil
}
