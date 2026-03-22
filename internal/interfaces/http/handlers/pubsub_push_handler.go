package handlers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/application/dto"
	"github.com/d-yanez/dys-sub-fsc-orders/internal/application/usecases"
)

type PubSubPushHandler struct {
	log     *slog.Logger
	useCase *usecases.ProcessOnOrderCreatedUseCase
}

func NewPubSubPushHandler(log *slog.Logger, useCase *usecases.ProcessOnOrderCreatedUseCase) *PubSubPushHandler {
	return &PubSubPushHandler{log: log, useCase: useCase}
}

func (h *PubSubPushHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "error": "method_not_allowed"})
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		h.respondJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "cannot_read_body"})
		return
	}

	var envelope dto.PubSubPushEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		h.respondJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "invalid_json"})
		return
	}

	if envelope.Message.Data == "" {
		h.log.Warn("pubsub_bad_envelope", "reason", "missing_message_data", "messageId", envelope.Message.MessageID)
		h.respondJSON(w, http.StatusOK, map[string]any{"ok": true, "status": "ignored", "reason": "missing_message_data"})
		return
	}

	raw, err := base64.StdEncoding.DecodeString(envelope.Message.Data)
	if err != nil {
		h.log.Warn("pubsub_bad_data", "reason", "invalid_base64", "messageId", envelope.Message.MessageID)
		h.respondJSON(w, http.StatusOK, map[string]any{"ok": true, "status": "ignored", "reason": "invalid_base64"})
		return
	}

	var event dto.FalabellaEvent
	if err := json.Unmarshal(raw, &event); err != nil {
		h.log.Warn("pubsub_bad_data", "reason", "invalid_event_json", "messageId", envelope.Message.MessageID)
		h.respondJSON(w, http.StatusOK, map[string]any{"ok": true, "status": "ignored", "reason": "invalid_event_json"})
		return
	}

	decision := h.useCase.Evaluate(event)
	h.log.Info("pubsub_event_received",
		"messageId", envelope.Message.MessageID,
		"eventType", decision.EventType,
		"orderId", decision.OrderID,
		"shouldProcess", decision.ShouldProcess,
		"reason", decision.Reason,
	)

	if !decision.ShouldProcess {
		h.respondJSON(w, http.StatusOK, map[string]any{
			"ok":     true,
			"status": "ignored",
			"reason": decision.Reason,
		})
		return
	}

	// PR2 scope: ingestion + parsing + routing decision only.
	h.respondJSON(w, http.StatusOK, map[string]any{
		"ok":        true,
		"status":    "accepted",
		"eventType": decision.EventType,
		"orderId":   decision.OrderID,
	})
}

func (h *PubSubPushHandler) respondJSON(w http.ResponseWriter, status int, body map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
