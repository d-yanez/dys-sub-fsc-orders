package usecases

import (
	"log/slog"
	"strings"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/application/dto"
)

type ProcessDecision struct {
	ShouldProcess bool
	EventType     string
	OrderID       string
	Reason        string
}

type ProcessOnOrderCreatedUseCase struct {
	log *slog.Logger
}

func NewProcessOnOrderCreatedUseCase(log *slog.Logger) *ProcessOnOrderCreatedUseCase {
	return &ProcessOnOrderCreatedUseCase{log: log}
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
