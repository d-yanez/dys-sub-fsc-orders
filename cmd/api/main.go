package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/application/usecases"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/infrastructure/fsc"
	mongoadapter "github.com/d-yanez/dys-sub-fsc-orders/internal/infrastructure/mongo"
	httpint "github.com/d-yanez/dys-sub-fsc-orders/internal/interfaces/http"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/interfaces/http/handlers"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/platform/config"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/platform/logger"
)

func main() {
	cfg := config.Load()
	appLogger := logger.New(cfg.ServiceName, cfg.LogLevel)
	slog.SetDefault(appLogger)

	ctx := context.Background()
	mongoClient, err := mongoadapter.NewClient(ctx, cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		appLogger.Error("mongo initialization failed", "error", err)
		return
	}
	defer func() {
		_ = mongoClient.Close(ctx)
	}()

	fscClient := fsc.NewClient(cfg.FSCBaseURL, cfg.FSCAPIKey, time.Duration(cfg.HTTPTimeoutMS)*time.Millisecond)
	orderRepo := mongoadapter.NewOrderRepository(mongoClient)
	orderItemRepo := mongoadapter.NewOrderItemRepository(mongoClient)
	eventLogRepo := mongoadapter.NewEventLogRepository(mongoClient)

	processEventUseCase := usecases.NewProcessOnOrderCreatedUseCase(appLogger, fscClient, orderRepo, orderItemRepo, eventLogRepo)
	pubSubHandler := handlers.NewPubSubPushHandler(appLogger, processEventUseCase)

	router := httpint.NewRouter(cfg, appLogger, pubSubHandler)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	appLogger.Info("starting dys-sub-fsc-orders", "addr", server.Addr, "environment", cfg.Environment)
	if err := server.ListenAndServe(); err != nil {
		appLogger.Error("server stopped with error", "error", err)
	}
}
