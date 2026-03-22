package usecases

import (
	"io"
	"log/slog"
	"testing"

	"github.com/d-yanez/dys-sub-fsc-orders/internal/application/dto"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestEvaluateAccepted(t *testing.T) {
	uc := NewProcessOnOrderCreatedUseCase(testLogger())
	evt := dto.FalabellaEvent{
		Event: "onOrderCreated",
	}
	evt.Payload.OrderID = "1146543495"

	got := uc.Evaluate(evt)
	if !got.ShouldProcess {
		t.Fatalf("expected ShouldProcess=true, got false reason=%s", got.Reason)
	}
}

func TestEvaluateIgnoredUnsupported(t *testing.T) {
	uc := NewProcessOnOrderCreatedUseCase(testLogger())
	evt := dto.FalabellaEvent{
		Event: "onOrderCanceled",
	}
	evt.Payload.OrderID = "1146543495"

	got := uc.Evaluate(evt)
	if got.ShouldProcess {
		t.Fatalf("expected ShouldProcess=false for unsupported event")
	}
	if got.Reason != "unsupported_event_type" {
		t.Fatalf("unexpected reason: %s", got.Reason)
	}
}
