package mongo

import (
	"context"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemRepository struct {
	client *Client
}

func NewOrderItemRepository(client *Client) *OrderItemRepository {
	return &OrderItemRepository{client: client}
}

func (r *OrderItemRepository) UpsertMany(ctx context.Context, items []entities.OrderItem) error {
	col := r.client.Collection("order_items")
	for _, item := range items {
		_, err := col.UpdateOne(ctx,
			bson.M{"_id": item.ID},
			bson.M{"$set": item},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
