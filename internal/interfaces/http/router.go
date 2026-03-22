package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/interfaces/http/handlers"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/interfaces/http/middleware"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/platform/config"
)

func NewRouter(cfg config.Config, log *slog.Logger, pubSubHandler *handlers.PubSubPushHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, "{\"status\":\"ok\",\"service\":\"%s\",\"environment\":\"%s\"}", cfg.ServiceName, cfg.Environment)
	})

	pushChain := middleware.OIDCStub(log, cfg.OIDCValidation, cfg.OIDCAudience, cfg.OIDCAllowedEmail, http.HandlerFunc(pubSubHandler.Handle))
	mux.Handle("POST /", pushChain)

	return middleware.RequestLogging(log, mux)
}
