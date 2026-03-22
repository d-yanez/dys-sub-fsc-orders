package handlers

import (
	"bytes"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/application/usecases"
)

func handlerTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestHandleAcceptedOnOrderCreated(t *testing.T) {
	uc := usecases.NewProcessOnOrderCreatedUseCase(handlerTestLogger())
	h := NewPubSubPushHandler(handlerTestLogger(), uc)

	eventJSON := `{"event":"onOrderCreated","payload":{"OrderId":"1146543495"}}`
	data := base64.StdEncoding.EncodeToString([]byte(eventJSON))
	body := `{"message":{"messageId":"m-1","data":"` + data + `"}}`

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"status":"accepted"`)) {
		t.Fatalf("expected accepted response, got body=%s", rec.Body.String())
	}
}

func TestHandleIgnoredMissingData(t *testing.T) {
	uc := usecases.NewProcessOnOrderCreatedUseCase(handlerTestLogger())
	h := NewPubSubPushHandler(handlerTestLogger(), uc)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"message":{"messageId":"m-1"}}`))
	rec := httptest.NewRecorder()
	h.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"status":"ignored"`)) {
		t.Fatalf("expected ignored response, got body=%s", rec.Body.String())
	}
}
