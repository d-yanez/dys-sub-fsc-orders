package ports

import (
	"context"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/entities"
)

type EventLogRepository interface {
	EnsureReceived(ctx context.Context, log entities.EventLog) error
	MarkProcessing(ctx context.Context, idempotencyKey, messageID, payloadHash, phase string) (bool, error)
	GetByID(ctx context.Context, idempotencyKey string) (*entities.EventLog, error)
	MarkCompleted(ctx context.Context, idempotencyKey string, status string, warning string) error
	MarkFailed(ctx context.Context, idempotencyKey string, phase string, errSummary string) error
}
