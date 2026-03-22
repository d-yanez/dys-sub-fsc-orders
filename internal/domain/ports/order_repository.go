package ports

import (
	"context"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/entities"
)

type OrderRepository interface {
	Upsert(ctx context.Context, order entities.Order) error
}
