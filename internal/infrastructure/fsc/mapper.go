package fsc

import (
	"strconv"
	"strings"
)

func firstString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		switch vv := v.(type) {
		case string:
			if vv != "" {
				return vv
			}
		case float64:
			return strconv.FormatInt(int64(vv), 10)
		case int:
			return strconv.Itoa(vv)
		}
	}
	return ""
}

func firstInt(m map[string]any, keys ...string) int {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		switch vv := v.(type) {
		case float64:
			return int(vv)
		case int:
			return vv
		case string:
			parsed, err := strconv.Atoi(vv)
			if err == nil {
				return parsed
			}
		}
	}
	return 0
}

func firstNumber(m map[string]any, keys ...string) *float64 {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		switch vv := v.(type) {
		case float64:
			n := vv
			return &n
		case float32:
			n := float64(vv)
			return &n
		case int:
			n := float64(vv)
			return &n
		case int64:
			n := float64(vv)
			return &n
		case string:
			raw := strings.TrimSpace(strings.ReplaceAll(vv, ",", ""))
			if raw == "" {
				continue
			}
			parsed, err := strconv.ParseFloat(raw, 64)
			if err == nil {
				n := parsed
				return &n
			}
		}
	}
	return nil
}

func firstInt64(m map[string]any, keys ...string) int64 {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		switch vv := v.(type) {
		case float64:
			return int64(vv)
		case float32:
			return int64(vv)
		case int:
			return int64(vv)
		case int64:
			return vv
		case string:
			raw := strings.TrimSpace(strings.ReplaceAll(vv, ",", ""))
			if raw == "" {
				continue
			}
			if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil {
				return parsed
			}
			if parsed, err := strconv.ParseFloat(raw, 64); err == nil {
				return int64(parsed)
			}
		}
	}
	return 0
}

func firstBool(m map[string]any, keys ...string) *bool {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		switch vv := v.(type) {
		case bool:
			b := vv
			return &b
		case string:
			raw := strings.ToLower(strings.TrimSpace(vv))
			if raw == "true" {
				b := true
				return &b
			}
			if raw == "false" {
				b := false
				return &b
			}
		}
	}
	return nil
}

func firstMap(m map[string]any, keys ...string) map[string]any {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		if mm, ok := v.(map[string]any); ok {
			return cloneMap(mm)
		}
	}
	return nil
}

func cloneMap(in map[string]any) map[string]any {
	if in == nil {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
