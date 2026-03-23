package dto

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type StringOrNumber string

func (s *StringOrNumber) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*s = ""
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		*s = StringOrNumber(asString)
		return nil
	}

	var asNumber json.Number
	if err := json.Unmarshal(data, &asNumber); err == nil {
		*s = StringOrNumber(asNumber.String())
		return nil
	}

	var asFloat float64
	if err := json.Unmarshal(data, &asFloat); err == nil {
		*s = StringOrNumber(strconv.FormatFloat(asFloat, 'f', -1, 64))
		return nil
	}

	return fmt.Errorf("unsupported OrderId value")
}

type FalabellaEvent struct {
	Event     string `json:"event"`
	EventType string `json:"eventType"`
	Payload   struct {
		OrderID StringOrNumber `json:"OrderId"`
	} `json:"payload"`
}

func (f FalabellaEvent) EffectiveEventType() string {
	if f.EventType != "" {
		return f.EventType
	}
	return f.Event
}
