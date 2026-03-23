package dto

import (
	"encoding/json"
	"testing"
)

func TestFalabellaEventUnmarshalOrderIDAsString(t *testing.T) {
	raw := []byte(`{"event":"onOrderCreated","payload":{"OrderId":"1147107464"}}`)
	var evt FalabellaEvent
	if err := json.Unmarshal(raw, &evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(evt.Payload.OrderID) != "1147107464" {
		t.Fatalf("unexpected order id: %s", evt.Payload.OrderID)
	}
}

func TestFalabellaEventUnmarshalOrderIDAsNumber(t *testing.T) {
	raw := []byte(`{"event":"onOrderCreated","payload":{"OrderId":1147107464}}`)
	var evt FalabellaEvent
	if err := json.Unmarshal(raw, &evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(evt.Payload.OrderID) != "1147107464" {
		t.Fatalf("unexpected order id: %s", evt.Payload.OrderID)
	}
}
