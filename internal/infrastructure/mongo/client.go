package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewClient(ctx context.Context, uri, dbName string) (*Client, error) {
	if uri == "" {
		return nil, errors.New("MONGODB_URI is required")
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := mongoClient.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	c := &Client{client: mongoClient, db: mongoClient.Database(dbName)}
	if err := c.ensureIndexes(ctx); err != nil {
		return nil, fmt.Errorf("ensure indexes: %w", err)
	}

	return c, nil
}

func (c *Client) Collection(name string) *mongo.Collection {
	return c.db.Collection(name)
}

func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}
