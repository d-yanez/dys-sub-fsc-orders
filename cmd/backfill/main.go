package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/infrastructure/fsc"
	mongoadapter "github.com/d-yanez/dys-sub-fsc-orders/internal/infrastructure/mongo"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/platform/config"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/platform/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderRow struct {
	ID string `bson:"_id"`
}

func main() {
	dryRun := flag.Bool("dry-run", true, "simulate updates without writing to Mongo")
	limit := flag.Int("limit", 500, "max orders to scan")
	batchSize := flag.Int("batch-size", 100, "Mongo cursor batch size")
	fromOrderID := flag.String("from-order-id", "", "resume scanning from this orderId (exclusive)")
	flag.Parse()

	cfg := config.Load()
	log := logger.New(cfg.ServiceName+"-backfill", cfg.LogLevel)
	slog.SetDefault(log)

	if cfg.MongoURI == "" {
		log.Error("MONGODB_URI is required")
		return
	}

	ctx := context.Background()
	mongoClient, err := mongoadapter.NewClient(ctx, cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Error("mongo initialization failed", "error", err)
		return
	}
	defer func() { _ = mongoClient.Close(ctx) }()

	fscClient := fsc.NewClient(cfg.FSCBaseURL, cfg.FSCAPIKey, time.Duration(cfg.HTTPTimeoutMS)*time.Millisecond)
	ordersCol := mongoClient.Collection("orders")

	filter := bson.M{
		"$or": []bson.M{
			{"financial": bson.M{"$exists": false}},
			{"financial.invoiceRequired": bson.M{"$exists": false}},
			{"financial.grandTotal": bson.M{"$exists": false}},
			{"financial.productTotal": bson.M{"$exists": false}},
			{"financial.taxAmount": bson.M{"$exists": false}},
			{"financial.shippingFeeTotal": bson.M{"$exists": false}},
			{"financial.documentType": bson.M{"$exists": false}},
			{"addresses": bson.M{"$exists": false}},
			{"addresses.billing": bson.M{"$exists": false}},
			{"addresses.shipping": bson.M{"$exists": false}},
		},
	}
	if *fromOrderID != "" {
		filter["_id"] = bson.M{"$gt": *fromOrderID}
	}

	opts := options.Find().
		SetProjection(bson.M{"_id": 1}).
		SetSort(bson.D{{Key: "_id", Value: 1}}).
		SetBatchSize(int32(*batchSize)).
		SetLimit(int64(*limit))

	cur, err := ordersCol.Find(ctx, filter, opts)
	if err != nil {
		log.Error("query missing orders failed", "error", err)
		return
	}
	defer cur.Close(ctx)

	var (
		scanned int
		updated int
		failed  int
	)

	for cur.Next(ctx) {
		var row orderRow
		if err := cur.Decode(&row); err != nil {
			failed++
			log.Error("decode row failed", "error", err)
			continue
		}
		scanned++

		orderResp, err := fscClient.GetOrder(ctx, row.ID)
		if err != nil {
			failed++
			log.Error("fetch order from dys-api-fsc failed", "orderId", row.ID, "error", err)
			continue
		}

		update := bson.M{
			"financial": bson.M{
				"grandTotal":       nullableFloat(orderResp.GrandTotal),
				"productTotal":     nullableFloat(orderResp.ProductTotal),
				"taxAmount":        nullableFloat(orderResp.TaxAmount),
				"shippingFeeTotal": nullableFloat(orderResp.ShippingFeeTotal),
				"invoiceRequired":  nullableBool(orderResp.InvoiceRequired),
				"documentType":     deriveDocumentType(orderResp.InvoiceRequired),
			},
			"addresses": bson.M{
				"billing":  nullableMap(orderResp.AddressBilling),
				"shipping": nullableMap(orderResp.AddressShipping),
			},
			"audit.updatedAt": time.Now().UTC(),
		}

		if *dryRun {
			log.Info("dry-run update candidate", "orderId", row.ID, "update", update)
			updated++
			continue
		}

		_, err = ordersCol.UpdateOne(ctx, bson.M{"_id": row.ID}, bson.M{"$set": update})
		if err != nil {
			failed++
			log.Error("update order failed", "orderId", row.ID, "error", err)
			continue
		}
		updated++
	}

	if err := cur.Err(); err != nil {
		log.Error("cursor error", "error", err)
	}

	log.Info("backfill finished",
		"dryRun", *dryRun,
		"scanned", scanned,
		"updated", updated,
		"failed", failed,
		"limit", *limit,
		"batchSize", *batchSize,
		"fromOrderID", *fromOrderID,
	)

	fmt.Printf("done scanned=%d updated=%d failed=%d dryRun=%v\n", scanned, updated, failed, *dryRun)
}

func nullableFloat(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}

func nullableBool(v *bool) any {
	if v == nil {
		return nil
	}
	return *v
}

func nullableMap(v map[string]any) any {
	if v == nil {
		return nil
	}
	return v
}

func deriveDocumentType(invoiceRequired *bool) any {
	if invoiceRequired == nil {
		return nil
	}
	if *invoiceRequired {
		return "FACTURA"
	}
	return "BOLETA"
}
