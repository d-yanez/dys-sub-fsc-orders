package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/application/usecases"
	httpint "github.com/d-yanez/dys-sub-fsc-orders/internal/interfaces/http"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/interfaces/http/handlers"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/platform/config"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/platform/logger"
)

func main() {
	cfg := config.Load()
	appLogger := logger.New(cfg.ServiceName, cfg.LogLevel)
	slog.SetDefault(appLogger)

	processEventUseCase := usecases.NewProcessOnOrderCreatedUseCase(appLogger)
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
