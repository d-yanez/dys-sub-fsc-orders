package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/application/dto"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/entities"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/enums"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/ports"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/services"
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
	eventLogRepo  ports.EventLogRepository
	telegram      ports.TelegramNotifier
	stockViewBase string
}

type ProcessResult struct {
	Status           enums.ResultStatus
	EventType        string
	OrderID          string
	OrderNumber      string
	ItemsCount       int
	Items            []ResultItem
	Warning          string
	Duplicate        bool
	Phase            string
	MessageID        string
	FirstOrderItemID string
	FirstSKU         string
	FirstItemName    string
	ThumbnailURL     string
	ErrorSummary     string
}

type ResultItem struct {
	OrderItemID string
	SKU         string
	Name        string
	Quantity    int
}

func NewProcessOnOrderCreatedUseCase(
	log *slog.Logger,
	fscClient ports.FSCClient,
	orderRepo ports.OrderRepository,
	orderItemRepo ports.OrderItemRepository,
	eventLogRepo ports.EventLogRepository,
	telegram ports.TelegramNotifier,
	stockViewBase string,
) *ProcessOnOrderCreatedUseCase {
	if strings.TrimSpace(stockViewBase) == "" {
		stockViewBase = "https://dy-api-utils-785293986978.us-central1.run.app/stock/view"
	}
	return &ProcessOnOrderCreatedUseCase{
		log:           log,
		fscClient:     fscClient,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		eventLogRepo:  eventLogRepo,
		telegram:      telegram,
		stockViewBase: stockViewBase,
	}
}

func (u *ProcessOnOrderCreatedUseCase) Evaluate(event dto.FalabellaEvent) ProcessDecision {
	eventType := strings.TrimSpace(event.EffectiveEventType())
	orderID := strings.TrimSpace(string(event.Payload.OrderID))

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

func (u *ProcessOnOrderCreatedUseCase) Process(ctx context.Context, event dto.FalabellaEvent, messageID, payloadHash string) (ProcessResult, error) {
	decision := u.Evaluate(event)
	if !decision.ShouldProcess {
		return ProcessResult{
			Status:    enums.ResultSuccess,
			EventType: decision.EventType,
			OrderID:   decision.OrderID,
			Warning:   decision.Reason,
			MessageID: messageID,
			Phase:     "IGNORED",
		}, nil
	}

	if u.fscClient == nil || u.orderRepo == nil || u.orderItemRepo == nil || u.eventLogRepo == nil {
		return ProcessResult{}, errors.New("usecase dependencies are not configured")
	}

	idempotencyKey := services.BuildIdempotencyKey(decision.EventType, decision.OrderID)
	now := time.Now().UTC()
	if err := u.eventLogRepo.EnsureReceived(ctx, entities.EventLog{
		ID:             idempotencyKey,
		IdempotencyKey: idempotencyKey,
		MessageID:      messageID,
		EventType:      decision.EventType,
		OrderID:        decision.OrderID,
		PayloadHash:    payloadHash,
		Status:         "RECEIVED",
		Processed:      false,
		Attempts:       0,
		Phase:          "RECEIVED",
		CreatedAt:      now,
		UpdatedAt:      now,
	}); err != nil {
		return ProcessResult{}, err
	}

	acquired, err := u.eventLogRepo.MarkProcessing(ctx, idempotencyKey, messageID, payloadHash, "PROCESSING")
	if err != nil {
		return ProcessResult{}, err
	}
	if !acquired {
		current, getErr := u.eventLogRepo.GetByID(ctx, idempotencyKey)
		if getErr != nil {
			return ProcessResult{}, getErr
		}
		if current != nil && current.Processed {
			u.log.Info("duplicate_event_ignored",
				"eventType", decision.EventType,
				"orderId", decision.OrderID,
				"messageId", messageID,
			)
			return ProcessResult{
				Status:    enums.ResultSuccess,
				EventType: decision.EventType,
				OrderID:   decision.OrderID,
				Warning:   "duplicate_event_ignored",
				Duplicate: true,
				MessageID: messageID,
				Phase:     "DUPLICATE",
			}, nil
		}
		return ProcessResult{
			EventType:    decision.EventType,
			OrderID:      decision.OrderID,
			MessageID:    messageID,
			Phase:        "PROCESSING",
			ErrorSummary: "event_in_processing_or_not_ready_for_retry",
		}, fmt.Errorf("event_in_processing_or_not_ready_for_retry")
	}

	baseResult := ProcessResult{
		EventType: decision.EventType,
		OrderID:   decision.OrderID,
		MessageID: messageID,
	}

	orderResp, err := u.fscClient.GetOrder(ctx, decision.OrderID)
	if err != nil {
		_ = u.eventLogRepo.MarkFailed(ctx, idempotencyKey, "FETCH_ORDER", err.Error())
		return u.fail(ctx, baseResult, "FETCH_ORDER", err)
	}
	itemsResp, err := u.fscClient.GetOrderItems(ctx, decision.OrderID)
	if err != nil {
		_ = u.eventLogRepo.MarkFailed(ctx, idempotencyKey, "FETCH_ORDER_ITEMS", err.Error())
		return u.fail(ctx, baseResult, "FETCH_ORDER_ITEMS", err)
	}

	now = time.Now().UTC()
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
		_ = u.eventLogRepo.MarkFailed(ctx, idempotencyKey, "PERSIST_ORDER", err.Error())
		baseResult.OrderNumber = order.OrderNumber
		return u.fail(ctx, baseResult, "PERSIST_ORDER", err)
	}

	partialWarning := ""
	orderItems := make([]entities.OrderItem, 0, len(itemsResp))
	resultItems := make([]ResultItem, 0, len(itemsResp))
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
		resultItems = append(resultItems, ResultItem{
			OrderItemID: itemID,
			SKU:         sku,
			Name:        strings.TrimSpace(it.Name),
			Quantity:    it.Quantity,
		})

		if baseResult.FirstOrderItemID == "" {
			baseResult.FirstOrderItemID = itemID
			baseResult.FirstSKU = sku
			baseResult.FirstItemName = strings.TrimSpace(it.Name)
			if thumbnail != nil {
				baseResult.ThumbnailURL = *thumbnail
			}
		}
	}

	if err := u.orderItemRepo.UpsertMany(ctx, orderItems); err != nil {
		_ = u.eventLogRepo.MarkFailed(ctx, idempotencyKey, "PERSIST_ORDER_ITEMS", err.Error())
		baseResult.OrderNumber = order.OrderNumber
		return u.fail(ctx, baseResult, "PERSIST_ORDER_ITEMS", err)
	}

	status := enums.ResultSuccess
	if partialWarning != "" {
		status = enums.ResultPartialSuccess
	}
	if err := u.eventLogRepo.MarkCompleted(ctx, idempotencyKey, string(status), partialWarning); err != nil {
		baseResult.OrderNumber = order.OrderNumber
		return u.fail(ctx, baseResult, "FINALIZE_EVENT_LOG", err)
	}

	result := ProcessResult{
		Status:           status,
		EventType:        decision.EventType,
		OrderID:          decision.OrderID,
		OrderNumber:      order.OrderNumber,
		ItemsCount:       len(orderItems),
		Items:            resultItems,
		Warning:          partialWarning,
		MessageID:        messageID,
		Phase:            "COMPLETED",
		FirstOrderItemID: baseResult.FirstOrderItemID,
		FirstSKU:         baseResult.FirstSKU,
		FirstItemName:    baseResult.FirstItemName,
		ThumbnailURL:     baseResult.ThumbnailURL,
	}
	u.notifyFinalNonBlocking(ctx, result, "")
	return result, nil
}

func (u *ProcessOnOrderCreatedUseCase) fail(ctx context.Context, base ProcessResult, phase string, err error) (ProcessResult, error) {
	result := base
	result.Status = enums.ResultFailed
	result.Phase = phase
	result.ErrorSummary = err.Error()
	u.notifyFinalNonBlocking(ctx, result, phase)
	return result, err
}

func (u *ProcessOnOrderCreatedUseCase) notifyFinalNonBlocking(ctx context.Context, result ProcessResult, failedPhase string) {
	if u.telegram == nil || result.Duplicate {
		return
	}
	msg := buildTelegramMessage(result, failedPhase, u.stockViewBase)
	if err := u.telegram.Send(ctx, msg); err != nil {
		u.log.Error("telegram_notification_failed",
			"orderId", result.OrderID,
			"eventType", result.EventType,
			"status", result.Status,
			"phase", result.Phase,
			"error", err.Error(),
		)
	}
}

func buildTelegramMessage(result ProcessResult, failedPhase string, stockViewBase string) ports.TelegramMessage {
	lines := make([]string, 0, 16)
	switch result.Status {
	case enums.ResultSuccess:
		lines = append(lines, "<b>Nueva orden Falabella procesada</b>")
	case enums.ResultPartialSuccess:
		lines = append(lines, "<b>Orden Falabella procesada con observaciones</b>")
	default:
		lines = append(lines, "<b>Orden Falabella no procesada</b>")
	}
	lines = append(lines,
		"eventType: "+htmlEsc(result.EventType),
		"orderId: <code>"+htmlEsc(fallback(result.OrderID, "N/A"))+"</code>",
		"orderNumber: <code>"+htmlEsc(fallback(result.OrderNumber, "N/A"))+"</code>",
		"itemsPersistidos: "+fmt.Sprintf("%d", result.ItemsCount),
		"orderItemId: <code>"+htmlEsc(fallback(result.FirstOrderItemID, "N/A"))+"</code>",
		"sku: <code>"+htmlEsc(fallback(result.FirstSKU, "N/A"))+"</code>",
		"resultado: "+htmlEsc(string(result.Status)),
		"messageId: <code>"+htmlEsc(fallback(result.MessageID, "N/A"))+"</code>",
		"phase: "+htmlEsc(fallback(result.Phase, "N/A")),
	)
	if result.FirstItemName != "" {
		lines = append(lines, "item: "+htmlEsc(result.FirstItemName))
	}
	if len(result.Items) > 0 {
		lines = append(lines, "items:")
		for idx, item := range result.Items {
			lines = append(lines,
				fmt.Sprintf("%d) orderItemId: <code>%s</code> | sku: <code>%s</code> | qty: <code>%d</code> | item: %s",
					idx+1,
					htmlEsc(fallback(item.OrderItemID, "N/A")),
					htmlEsc(fallback(item.SKU, "N/A")),
					item.Quantity,
					htmlEsc(fallback(item.Name, "(sin nombre)")),
				),
			)
		}
	}
	if failedPhase != "" {
		lines = append(lines, "failedPhase: "+htmlEsc(failedPhase))
	}
	if result.Warning != "" {
		lines = append(lines, "warning: "+htmlEsc(result.Warning))
	}
	if result.ErrorSummary != "" {
		lines = append(lines, "error: "+htmlEsc(result.ErrorSummary))
	}
	if stockViewURL := buildStockViewURL(stockViewBase, result.FirstSKU); stockViewURL != "" {
		lines = append(lines, `ver stock bodega: <a href="`+htmlEsc(stockViewURL)+`">Ver stock bodega</a>`)
	}
	lines = append(lines, "timestamp: "+time.Now().Format(time.RFC3339))
	return ports.TelegramMessage{
		Text:           strings.Join(lines, "\n"),
		PhotoURL:       result.ThumbnailURL,
		ParseMode:      "HTML",
		DisablePreview: true,
	}
}

func buildStockViewURL(base, sku string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	sku = strings.TrimSpace(sku)
	if base == "" || sku == "" {
		return ""
	}
	return base + "/" + url.PathEscape(sku)
}

func fallback(v string, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}

func htmlEsc(v string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
	)
	return replacer.Replace(v)
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
