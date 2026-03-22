package middleware

import (
	"log/slog"
	"net/http"
)

func OIDCStub(logger *slog.Logger, enabled bool, audience, allowedEmail string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if enabled {
			logger.Debug("oidc_validation_stub_enabled",
				"audience", audience,
				"allowed_email", allowedEmail,
			)
		}
		next.ServeHTTP(w, r)
	})
}
