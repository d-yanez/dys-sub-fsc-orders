package mongo

import (
	"context"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
	client *Client
}

func NewOrderRepository(client *Client) *OrderRepository {
	return &OrderRepository{client: client}
}

func (r *OrderRepository) Upsert(ctx context.Context, order entities.Order) error {
	col := r.client.Collection("orders")
	_, err := col.UpdateOne(ctx,
		bson.M{"_id": order.ID},
		bson.M{"$set": order},
		options.Update().SetUpsert(true),
	)
	return err
}
