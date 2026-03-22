package ports

import (
	"context"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/entities"
)

type OrderItemRepository interface {
	UpsertMany(ctx context.Context, items []entities.OrderItem) error
}
