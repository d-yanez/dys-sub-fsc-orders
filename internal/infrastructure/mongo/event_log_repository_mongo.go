package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventLogRepository struct {
	client *Client
}

func NewEventLogRepository(client *Client) *EventLogRepository {
	return &EventLogRepository{client: client}
}

func (r *EventLogRepository) EnsureReceived(ctx context.Context, log entities.EventLog) error {
	col := r.client.Collection("event_logs")
	now := time.Now().UTC()
	_, err := col.UpdateOne(
		ctx,
		bson.M{"_id": log.ID},
		bson.M{
			"$setOnInsert": bson.M{
				"_id":            log.ID,
				"idempotencyKey": log.IdempotencyKey,
				"eventType":      log.EventType,
				"orderId":        log.OrderID,
				"status":         "RECEIVED",
				"processed":      false,
				"attempts":       0,
				"phase":          "RECEIVED",
				"createdAt":      now,
			},
			"$set": bson.M{
				"updatedAt": now,
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *EventLogRepository) MarkProcessing(ctx context.Context, idempotencyKey, messageID, payloadHash, phase string) (bool, error) {
	col := r.client.Collection("event_logs")
	now := time.Now().UTC()
	res, err := col.UpdateOne(
		ctx,
		bson.M{
			"_id":       idempotencyKey,
			"processed": false,
			"status": bson.M{
				"$in": []string{"RECEIVED", "FAILED"},
			},
		},
		bson.M{
			"$set": bson.M{
				"status":      "PROCESSING",
				"phase":       phase,
				"messageId":   messageID,
				"payloadHash": payloadHash,
				"updatedAt":   now,
			},
			"$inc": bson.M{"attempts": 1},
		},
	)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount == 1, nil
}

func (r *EventLogRepository) GetByID(ctx context.Context, idempotencyKey string) (*entities.EventLog, error) {
	col := r.client.Collection("event_logs")
	var out entities.EventLog
	err := col.FindOne(ctx, bson.M{"_id": idempotencyKey}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func (r *EventLogRepository) MarkCompleted(ctx context.Context, idempotencyKey string, status string, warning string) error {
	col := r.client.Collection("event_logs")
	now := time.Now().UTC()
	update := bson.M{
		"$set": bson.M{
			"status":      status,
			"processed":   true,
			"phase":       "COMPLETED",
			"processedAt": now,
			"updatedAt":   now,
		},
	}
	if warning != "" {
		update["$set"].(bson.M)["errorSummary"] = warning
	}
	_, err := col.UpdateOne(ctx, bson.M{"_id": idempotencyKey}, update)
	return err
}

func (r *EventLogRepository) MarkFailed(ctx context.Context, idempotencyKey string, phase string, errSummary string) error {
	col := r.client.Collection("event_logs")
	now := time.Now().UTC()
	_, err := col.UpdateOne(
		ctx,
		bson.M{"_id": idempotencyKey},
		bson.M{
			"$set": bson.M{
				"status":       "FAILED",
				"processed":    false,
				"phase":        phase,
				"errorSummary": errSummary,
				"updatedAt":    now,
			},
		},
	)
	return err
}
