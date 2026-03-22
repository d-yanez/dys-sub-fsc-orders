package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *Client) ensureIndexes(ctx context.Context) error {
	orders := c.Collection("orders")
	_, err := orders.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "orderNumber", Value: 1}},
			Options: options.Index().SetName("ux_orders_orderNumber").SetUnique(true).SetSparse(true),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}, {Key: "createdAt", Value: -1}},
			Options: options.Index().SetName("ix_orders_status_createdAt"),
		},
	})
	if err != nil {
		return err
	}

	orderItems := c.Collection("order_items")
	_, err = orderItems.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "orderId", Value: 1}},
			Options: options.Index().SetName("ix_order_items_orderId"),
		},
		{
			Keys:    bson.D{{Key: "orderId", Value: 1}, {Key: "sku", Value: 1}},
			Options: options.Index().SetName("ix_order_items_orderId_sku"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("ix_order_items_status"),
		},
	})
	if err != nil {
		return err
	}

	eventLogs := c.Collection("event_logs")
	_, err = eventLogs.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "messageId", Value: 1}},
			Options: options.Index().SetName("ix_event_logs_messageId"),
		},
		{
			Keys:    bson.D{{Key: "orderId", Value: 1}, {Key: "updatedAt", Value: -1}},
			Options: options.Index().SetName("ix_event_logs_orderId_updatedAt"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}, {Key: "updatedAt", Value: -1}},
			Options: options.Index().SetName("ix_event_logs_status_updatedAt"),
		},
	})

	return err
}
